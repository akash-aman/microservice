package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"pkg/kafka"
	"pkg/logger"
	"sort"
	"time"
	kfka "github.com/segmentio/kafka-go"
)

const (
	UserEventTopic       = "user-events"
	UserRecommendations  = "user-recommendations"
	BatchProcessInterval = 1 * time.Hour
)

// Event types
const (
	EventTypeView     = "view"
	EventTypePurchase = "purchase"
	EventTypeRate     = "rate"
	EventTypeSearch   = "search"
)

// UserEvent represents an action taken by a user
type UserEvent struct {
	UserID    string                 `json:"userId"`
	ItemID    string                 `json:"itemId"`
	EventType string                 `json:"eventType"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Recommendation represents a recommended item for a specific user
type Recommendation struct {
	UserID      string    `json:"userId"`
	ItemID      string    `json:"itemId"`
	Score       float64   `json:"score"`
	Reason      string    `json:"reason"`
	GeneratedAt time.Time `json:"generatedAt"`
}

// RecommendationService manages the recommendation operations
type RecommendationService struct {
	kafkaClient *kafka.KafkaClient
	log         logger.Zapper
	// In-memory store for simple implementation
	// In production, you'd use a database or cache
	userProfiles     map[string]UserProfile
	itemProfiles     map[string]ItemProfile
	recommendationDB map[string][]Recommendation
}

// UserProfile maintains user behavior data
type UserProfile struct {
	UserID         string
	ViewedItems    map[string]int
	PurchasedItems map[string]int
	RatedItems     map[string]float64
	Searches       []string
	LastUpdated    time.Time
}

// ItemProfile maintains item metadata and statistics
type ItemProfile struct {
	ItemID        string
	Category      string
	Tags          []string
	ViewCount     int
	PurchaseCount int
	AvgRating     float64
	LastUpdated   time.Time
}

// NewRecommendationService creates a new recommendation service
func NewRecommendationService(kafkaClient *kafka.KafkaClient, log logger.Zapper) *RecommendationService {
	return &RecommendationService{
		kafkaClient:      kafkaClient,
		log:              log,
		userProfiles:     make(map[string]UserProfile),
		itemProfiles:     make(map[string]ItemProfile),
		recommendationDB: make(map[string][]Recommendation),
	}
}

// Start begins the recommendation service processes
func (s *RecommendationService) Start(ctx context.Context) error {
	// Start the event consumer
	go func() {
		err := s.consumeUserEvents(ctx)
		if err != nil {
			s.log.Errorf(ctx, "Failed to consume user events: %v", err)
		}
	}()

	// Start the batch recommendation generator
	go s.scheduleBatchRecommendations(ctx)

	return nil
}

// TrackUserEvent tracks a user interaction with an item
func (s *RecommendationService) TrackUserEvent(ctx context.Context, event UserEvent) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal user event: %w", err)
	}

	producer, err := s.kafkaClient.GetProducer(ctx, UserEventTopic, kafka.KafkaConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   UserEventTopic,
	})
	if err != nil {
		return fmt.Errorf("failed to get Kafka producer: %w", err)
	}

	key := []byte(event.UserID)
	err = producer.Produce(ctx, key, eventData)
	if err != nil {
		return fmt.Errorf("failed to produce user event: %w", err)
	}

	return nil
}

// GetRecommendations retrieves recommendations for a user
func (s *RecommendationService) GetRecommendations(ctx context.Context, userID string, limit int) ([]Recommendation, error) {
	recommendations, exists := s.recommendationDB[userID]
	if !exists || len(recommendations) == 0 {
		// Generate real-time recommendations if none exist
		var err error
		recommendations, err = s.generateRecommendationsForUser(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate recommendations: %w", err)
		}
	}

	// Apply limit
	if limit > 0 && len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

// consumeUserEvents processes incoming user events
func (s *RecommendationService) consumeUserEvents(ctx context.Context) error {
	consumer, err := s.kafkaClient.GetConsumer(ctx, UserEventTopic, kafka.KafkaConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   UserEventTopic,
		GroupID: "recommendation-processor",
	})
	if err != nil {
		return fmt.Errorf("failed to get Kafka consumer: %w", err)
	}

	handler := func(data []byte) error {
		var event UserEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("failed to unmarshal user event: %w", err)
		}

		// Update user profile
		s.updateUserProfile(ctx, event)

		// Update item profile
		s.updateItemProfile(ctx, event)

		return nil
	}

	return consumer.ConsumeWithCommit(ctx, handler)
}

// updateUserProfile updates the user's profile based on their activity
func (s *RecommendationService) updateUserProfile(ctx context.Context, event UserEvent) {
	profile, exists := s.userProfiles[event.UserID]
	if !exists {
		profile = UserProfile{
			UserID:         event.UserID,
			ViewedItems:    make(map[string]int),
			PurchasedItems: make(map[string]int),
			RatedItems:     make(map[string]float64),
			Searches:       []string{},
		}
	}

	// Update profile based on event type
	switch event.EventType {
	case EventTypeView:
		profile.ViewedItems[event.ItemID]++
	case EventTypePurchase:
		profile.PurchasedItems[event.ItemID]++
	case EventTypeRate:
		if rating, ok := event.Metadata["rating"].(float64); ok {
			profile.RatedItems[event.ItemID] = rating
		}
	case EventTypeSearch:
		if query, ok := event.Metadata["query"].(string); ok {
			profile.Searches = append(profile.Searches, query)
			// Keep only the last 100 searches
			if len(profile.Searches) > 100 {
				profile.Searches = profile.Searches[len(profile.Searches)-100:]
			}
		}
	}

	profile.LastUpdated = time.Now()
	s.userProfiles[event.UserID] = profile
}

// updateItemProfile updates the item's profile based on user activities
func (s *RecommendationService) updateItemProfile(ctx context.Context, event UserEvent) {
	profile, exists := s.itemProfiles[event.ItemID]
	if !exists {
		profile = ItemProfile{
			ItemID: event.ItemID,
		}

		// Extract item metadata if available
		if category, ok := event.Metadata["category"].(string); ok {
			profile.Category = category
		}
		if tags, ok := event.Metadata["tags"].([]string); ok {
			profile.Tags = tags
		}
	}

	// Update profile based on event type
	switch event.EventType {
	case EventTypeView:
		profile.ViewCount++
	case EventTypePurchase:
		profile.PurchaseCount++
	case EventTypeRate:
		if rating, ok := event.Metadata["rating"].(float64); ok {
			// Simple running average
			oldTotal := profile.AvgRating * float64(profile.ViewCount)
			newTotal := oldTotal + rating
			profile.AvgRating = newTotal / float64(profile.ViewCount+1)
		}
	}

	profile.LastUpdated = time.Now()
	s.itemProfiles[event.ItemID] = profile
}

// scheduleBatchRecommendations periodically generates recommendations for all users
func (s *RecommendationService) scheduleBatchRecommendations(ctx context.Context) {
	ticker := time.NewTicker(BatchProcessInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.generateBatchRecommendations(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// generateBatchRecommendations generates recommendations for all users
func (s *RecommendationService) generateBatchRecommendations(ctx context.Context) {
	s.log.Infof(ctx, "Starting batch recommendation generation")

	// Process each user
	for userID := range s.userProfiles {
		recommendations, err := s.generateRecommendationsForUser(ctx, userID)
		if err != nil {
			s.log.Errorf(ctx, "Failed to generate recommendations for user %s: %v", userID, err)
			continue
		}

		// Store recommendations
		s.recommendationDB[userID] = recommendations

		// Publish recommendations to Kafka
		s.publishRecommendations(ctx, userID, recommendations)
	}

	s.log.Infof(ctx, "Completed batch recommendation generation")
}

// generateRecommendationsForUser creates personalized recommendations for a specific user
func (s *RecommendationService) generateRecommendationsForUser(ctx context.Context, userID string) ([]Recommendation, error) {
	userProfile, exists := s.userProfiles[userID]
	if !exists {
		return []Recommendation{}, nil
	}

	// A simple scoring algorithm based on collaborative filtering
	// In production, you'd use more sophisticated algorithms (matrix factorization, neural networks, etc.)

	// Map to store item scores
	itemScores := make(map[string]float64)
	itemReasons := make(map[string]string)

	// 1. Content-based recommendations from viewed items
	for viewedItemID, viewCount := range userProfile.ViewedItems {
		viewedItem, exists := s.itemProfiles[viewedItemID]
		if !exists {
			continue
		}

		// Find similar items by category and tags
		for candidateItemID, candidateItem := range s.itemProfiles {
			if candidateItemID == viewedItemID {
				continue // Skip the same item
			}

			var score float64 = 0
			reason := ""

			// Same category bonus
			if viewedItem.Category == candidateItem.Category && viewedItem.Category != "" {
				score += 2.0
				reason = fmt.Sprintf("Similar to %s you viewed", viewedItemID)
			}

			// Tag overlap
			commonTags := 0
			for _, tag1 := range viewedItem.Tags {
				for _, tag2 := range candidateItem.Tags {
					if tag1 == tag2 {
						commonTags++
					}
				}
			}
			if commonTags > 0 {
				score += float64(commonTags) * 0.5
				if reason == "" {
					reason = fmt.Sprintf("Has tags similar to %s you viewed", viewedItemID)
				}
			}

			// Weight by view count
			score *= float64(viewCount)

			// Add to scores
			if score > 0 {
				itemScores[candidateItemID] += score
				if itemReasons[candidateItemID] == "" || len(itemReasons[candidateItemID]) > len(reason) {
					itemReasons[candidateItemID] = reason
				}
			}
		}
	}

	// 2. Boost scores for items purchased by similar users
	// Simple collaborative filtering approach
	for otherUserID, otherProfile := range s.userProfiles {
		if otherUserID == userID {
			continue
		}

		// Calculate user similarity based on common viewed items
		commonViews := 0
		for itemID := range userProfile.ViewedItems {
			if _, viewed := otherProfile.ViewedItems[itemID]; viewed {
				commonViews++
			}
		}

		similarity := float64(commonViews) / float64(len(userProfile.ViewedItems)+len(otherProfile.ViewedItems)-commonViews)
		if similarity < 0.1 {
			continue // Skip users with low similarity
		}

		// Add purchased items from similar users
		for purchasedItemID := range otherProfile.PurchasedItems {
			if _, alreadyViewed := userProfile.ViewedItems[purchasedItemID]; !alreadyViewed {
				itemScores[purchasedItemID] += similarity * 3.0
				if itemReasons[purchasedItemID] == "" {
					itemReasons[purchasedItemID] = "Popular with users who have similar tastes"
				}
			}
		}
	}

	// 3. Boost highly rated items
	for itemID, profile := range s.itemProfiles {
		if profile.AvgRating > 4.0 && profile.ViewCount > 10 {
			itemScores[itemID] += profile.AvgRating * 0.5
			if itemReasons[itemID] == "" {
				itemReasons[itemID] = "Highly rated by other users"
			}
		}
	}

	// Convert map to slice for sorting
	var recommendations []Recommendation
	now := time.Now()

	for itemID, score := range itemScores {
		// Skip items the user has already purchased
		if _, purchased := userProfile.PurchasedItems[itemID]; purchased {
			continue
		}

		recommendations = append(recommendations, Recommendation{
			UserID:      userID,
			ItemID:      itemID,
			Score:       score,
			Reason:      itemReasons[itemID],
			GeneratedAt: now,
		})
	}

	// Sort recommendations by score (descending)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Limit to top 50 recommendations
	if len(recommendations) > 50 {
		recommendations = recommendations[:50]
	}

	return recommendations, nil
}

// publishRecommendations publishes user recommendations to Kafka
func (s *RecommendationService) publishRecommendations(ctx context.Context, userID string, recommendations []Recommendation) error {
	if len(recommendations) == 0 {
		return nil
	}

	producer, err := s.kafkaClient.GetProducer(ctx, UserRecommendations, kafka.KafkaConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   UserRecommendations,
	})
	if err != nil {
		return fmt.Errorf("failed to get Kafka producer: %w", err)
	}

	// Create a message for each recommendation
	var messages []kfka.Message
	for _, rec := range recommendations {
		data, err := json.Marshal(rec)
		if err != nil {
			s.log.Errorf(ctx, "Failed to marshal recommendation: %v", err)
			continue
		}

		messages = append(messages, kfka.Message{
			Key:   []byte(rec.UserID),
			Value: data,
			Time:  time.Now(),
		})
	}

	// Send batch of recommendations
	if len(messages) > 0 {
		err = producer.ProduceBatch(ctx, messages)
		if err != nil {
			return fmt.Errorf("failed to produce recommendations: %w", err)
		}
	}

	return nil
}

// InitializeFromHistory loads historical data to bootstrap the recommendation system
func (s *RecommendationService) InitializeFromHistory(ctx context.Context, userEvents []UserEvent, itemMetadata map[string]map[string]interface{})  {
	// Process user events to build profiles
	for _, event := range userEvents {
		s.updateUserProfile(ctx, event)
		s.updateItemProfile(ctx, event)
	}

	// Enrich item profiles with metadata
	for itemID, metadata := range itemMetadata {
		profile, exists := s.itemProfiles[itemID]
		if !exists {
			profile = ItemProfile{
				ItemID: itemID,
			}
		}

		if category, ok := metadata["category"].(string); ok {
			profile.Category = category
		}
		if tags, ok := metadata["tags"].([]string); ok {
			profile.Tags = tags
		}

		s.itemProfiles[itemID] = profile
	}

	// Generate initial recommendations
	s.generateBatchRecommendations(ctx)
}
