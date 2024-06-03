#include <iostream>
#include <cuda_runtime.h>

__global__ void fibonacciGPU(int n, unsigned long long* result)
{
    int tid = threadIdx.x;

    if (tid == 0)
    {
        result[0] = 0; // First Fibonacci number
        result[1] = 1; // Second Fibonacci number

        for (int i = 2; i <= n; i++)
        {
            result[i] = result[i - 1] + result[i - 2];
        }
    }
}

int main()
{
    int n = 10; // Number of Fibonacci numbers to compute

    // Allocate memory on the GPU
    unsigned long long* d_result;
    cudaMalloc((void**)&d_result, (n + 1) * sizeof(unsigned long long));

    // Launch the kernel on the GPU
    fibonacciGPU<<<1, 1>>>(n, d_result);

    // Copy the result back to the host
    unsigned long long* h_result = new unsigned long long[n + 1];
    cudaMemcpy(h_result, d_result, (n + 1) * sizeof(unsigned long long), cudaMemcpyDeviceToHost);

    // Print the Fibonacci numbers
    for (int i = 0; i <= n; i++)
    {
        std::cout << h_result[i] << " ";
    }
    std::cout << std::endl;

    // Cleanup
    delete[] h_result;
    cudaFree(d_result);

    return 0;
}
