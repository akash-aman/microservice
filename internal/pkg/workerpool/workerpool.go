package workerpool

// TO DOs :
// - OpenTelemetry integration: Add hooks for exporting metrics to monitoring systems. Enhanced Control Features.
// - Dynamic configuration: Allow changing pool parameters (min/max workers, idle timeout) at runtime.
// - Result collection: Add support for tasks that return results and errors.
// - Task retry mechanism: Automatically retry failed tasks with configurable backoff.
// - Priority queue system: Implement multiple priority levels for tasks, where high-priority tasks are processed before lower-priority ones.
// - Circuit breaker pattern: Detect when too many tasks are failing and temporarily stop accepting new tasks for x time. -> 25% -> 50% -> 75% -> 100% if not failing. if failing then again open state.
// - Named worker groups: Create specialized worker groups for different types of tasks (e.g., IO-bound vs CPU-bound tasks).
// - Result collection: Add support for tasks that return results and errors.
// - Worker middleware: Allow registering middleware functions that wrap task execution (for logging, timing, etc.).
// - Task decorators: Provide helper functions to wrap tasks with common behaviors like retries, timeouts, etc.
// - Task timeout handling: Imagine you have a worker pool that processes various tasks, some of which might involve network calls, database queries, or other operations that could potentially hang or take too long. By implementing task timeout handling, you can ensure that each task is given a specific amount of time to complete. If a task exceeds this time, it is canceled, and the worker can move on to the next task.
// - Graceful shutdown with timeout: Imagine you have a worker pool processing various tasks, and you need to shut down the application. You want to ensure that the worker pool stops accepting new tasks and gives the currently running tasks a chance to complete. However, you also want to avoid waiting indefinitely for tasks to finish, so you specify a timeout for the shutdown process. If the timeout is reached, any remaining tasks are canceled.
// - Resource-aware scheduling: Assign tasks to workers based on available system resources like CPU and memory.
// - Persistent task queue: Store pending tasks to disk or database to survive application restarts.
// - Dead letter queue: Move failed tasks that have exceeded retry limits to a separate queue for inspection.
// - Distributed worker pools: Support worker pools across multiple machines or containers.
// - Batching capability: Group similar tasks together for more efficient processing.
// - Task scheduling: Allow tasks to be scheduled for future execution at specific times.
// - Adaptive scaling: Automatically adjust pool size.
