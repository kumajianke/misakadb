import time

from apis import api

HOST = "127.0.0.1"
PORT = 10032
TOTAL_REQUESTS = 1000000
REQUESTS_PER_CONNECTION = 5000


def run_batch(batch_requests: int) -> tuple[int, float]:
    start = time.perf_counter()
    success_count = 0

    with api.MisakaDBClient(HOST, PORT) as client:
        for _ in range(batch_requests):
            result = client.execute_command("get-service-info")
            if result is not None:
                success_count += 1

    elapsed = time.perf_counter() - start
    return success_count, elapsed


def main() -> None:
    success_times = 0
    batch_count = (TOTAL_REQUESTS + REQUESTS_PER_CONNECTION - 1) // REQUESTS_PER_CONNECTION
    total_start = time.perf_counter()

    for batch_index in range(batch_count):
        remaining = TOTAL_REQUESTS - batch_index * REQUESTS_PER_CONNECTION
        current_batch_requests = min(REQUESTS_PER_CONNECTION, remaining)

        print(
            f"第 {batch_index + 1}/{batch_count} 批压测: "
            f"单连接请求数 = {current_batch_requests}"
        )

        batch_success, batch_elapsed = run_batch(current_batch_requests)
        success_times += batch_success

        print(
            f"本批成功: {batch_success}/{current_batch_requests}, "
            f"耗时: {batch_elapsed:.4f} 秒"
        )

    total_elapsed = time.perf_counter() - total_start
    qps = success_times / total_elapsed if total_elapsed > 0 else 0.0
    success_rate = (success_times / TOTAL_REQUESTS * 100) if TOTAL_REQUESTS > 0 else 0.0

    print(f"总请求数: {TOTAL_REQUESTS}")
    print(f"成功次数: {success_times}")
    print(f"总耗时: {total_elapsed:.4f} 秒")
    print(f"响应率: {success_rate:.2f}%")
    print(f"平均 QPS: {qps:.2f}")


main()
