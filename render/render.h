#include <stdbool.h>

#ifndef RENDER_H
#define RENDER_H

typedef struct Image {
    void *pixel;
    int width;
    int height;
} Image;

bool Render(
        float angle, float tilt, float zoom,
        bool shadow, bool lighting,
        bool portrait, bool full,
        bool overlay, bool newSkin,
        Image result,
        Image head, Image *headOverlay,
        Image *body, Image *bodyOverlay,
        Image *leftArm, Image *leftArmOverlay, Image *rightArm, Image *rightArmOverlay,
        Image *leftLeg, Image *leftLegOverlay, Image *rightLeg, Image *rightLegOverlay);

#endif
