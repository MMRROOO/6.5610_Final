#include"util.h"
// mman library to be used for hugepage allocations (e.g. mmap or posix_memalign only)
#include <sys/mman.h>

#define CACHE_LINE_SIZE 64
#define CACHE_SETS 24576 
#define CACHE_WAYS 12


int main(){
    // create large page to work with more cache sizes
    void *buf= mmap(NULL, BUFSIZ, PROT_READ | PROT_WRITE, 
                MAP_POPULATE | MAP_ANONYMOUS | MAP_PRIVATE | MAP_HUGETLB, 
                -1, 0);
    if (buf == (void*) -1) {
	perror("mmap() error\n");
	exit(EXIT_FAILURE);
    }
    *((char *)buf) = 1;
    while (true){//evict cache set from LLC cache
	int SET = 0;
	for(int w=0; w<CACHE_WAYS; w++){
	    char tmp = *((char*)buf + (SET*CACHE_LINE_SIZE)+(w*CACHE_SETS*CACHE_LINE_SIZE));
	}
    }	
}
