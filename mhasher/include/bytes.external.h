
#pragma once

#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif
	void Sort(char* bytes, uint32_t* lengthVec, uint32_t count, uint32_t* indices);
	void Unique(char* bytes, uint32_t* lengthVec, uint32_t count, uint8_t* newIndices);
	void Remove(char* bytes, uint32_t* lengths, uint32_t count, char* toRemove, uint32_t* toRemoveLengths, uint32_t toRemoveCount, uint8_t* newIndices); //Remove only
	
	void* Start();
	void ToBuffer(void* engine, char* path, char* pathLen, uint32_t count);
	void FromBuffer(void* engine, char* buffer, char* outputLens, uint32_t* count);
	void Clear(void* engine);
	void Stop(void* engine);
#ifdef __cplusplus
}
#endif
