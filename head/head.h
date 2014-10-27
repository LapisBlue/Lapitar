#include <stdbool.h>

#ifndef HEAD_H
#define HEAD_H

bool RenderHead(
        float angle,
        bool shadow, bool lighting,
        void *result, int width, int height,
        void *head, int headWidth, int headHeight,
        void *helm, int helmWidth, int helmHeight);

#endif
