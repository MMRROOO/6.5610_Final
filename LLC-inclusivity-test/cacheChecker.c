#include"util.h"
// mman library to be used for hugepage allocations (e.g. mmap or posix_memalign only)
#include <sys/mman.h>

#define CACHE_LINE_SIZE 64
#define CACHE_SETS 24576 
#define TRIALS 10
#define LLCITRS 100000
#define ITRS 100000
#define L1_LINE_SIZE 64
#define L1_WAYS 12
#define L1_SETS 64
#define L2_LINE_SIZE 64 
#define L2_WAYS 10
#define L2_SETS 2048

int main() {
    uint8_t* access_times_LLC = malloc(sizeof(int)*TRIALS);
    
    int SET = 0;
    int sum = 0;
    int WAY = 0;

    void *buf= mmap(NULL, BUFSIZ, PROT_READ | PROT_WRITE, 
                MAP_POPULATE | MAP_ANONYMOUS | MAP_PRIVATE | MAP_HUGETLB, 
                -1, 0);
    if (buf == (void*) -1) {
		perror("mmap() error\n");
		exit(EXIT_FAILURE);
    }

    volatile uint8_t *target_buffer = (uint8_t *)malloc(CACHE_LINE_SIZE);
	int randomOffset[10] = {1, 32, 8, 9, 87, 11, 100, 93, 56, 28};

    for(int t=0; t<TRIALS; t++){
    	for (int i=0; i<LLCITRS; i++){
	// put target in l1 cache
			uint8_t tmp = target_buffer[randomOffset[i%10]];

			// volatile uint8_t *buff3 = (uint8_t*)malloc(L1_SIZE sizeof(uint8_t));
			
			for (int j=0; j<L1_WAYS; j++){ // fill set
				char tmp = *((char*)(buf + randomOffset[i%10]+ (j*L1_SETS*L1_LINE_SIZE))); //move to correct set then iterate to fill set
				tmp +=1;
			}
			// free(buff3);


			// volatile uint8_t *buff4 = (uint8_t*)malloc(8*L2_SIZE*sizeof(uint8_t));
			for (int j=0; j<L2_WAYS; j++){ // fill set
				char tmp = *((char*)(buf +randomOffset[i%10]+ (j*L2_SETS*L2_LINE_SIZE))); //move to correct set then iterate to fill set
				tmp +=1;
			}
			// free(buff4);

			clflush(((char*)buf + randomOffset[i%10]));
			//measure latency
			sum += measure_one_block_access_time(((char*)buf + randomOffset[i%10]));
        }	

		//printf("done trial: %d\n", t);

		access_times_LLC[t] = sum/LLCITRS;
		sum = 0;
    }

	free(target_buffer);

	printf("done checking l3 times\n");

	for (int t=0; t<TRIALS; t++){
		printf("LLC time: %d\n", access_times_LLC[t]);
    }
    
    uint8_t* access_times = alloca(sizeof(int)*TRIALS);
    


    *((char *)buf) = 1;
    for(int t=0; t<TRIALS; t++){
    	for(int i=0; i<ITRS; i++){
			char tmp2 = *(char*)buf + (SET*CACHE_LINE_SIZE)+(WAY*CACHE_SETS*CACHE_LINE_SIZE);

            int tmp = measure_one_block_access_time((char*)buf + (SET*CACHE_LINE_SIZE)+(WAY*CACHE_SETS*CACHE_LINE_SIZE));
	    	sum+= tmp;
		}
		access_times[t] = sum/ITRS;
		sum = 0;
    }

    for (int t=0; t<TRIALS; t++){
		printf("NormalAccess: %d\n", access_times[t]);
    }
    

}
