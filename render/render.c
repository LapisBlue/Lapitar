#include <stdio.h>
#include <GL/osmesa.h>
#include <GL/glu.h>
#include "render.h"
#include "render_constants.h"

OSMesaContext initOSMesa(Image result) {
    OSMesaContext ctx = OSMesaCreateContextExt(OSMESA_RGBA, 16, 0, 0, NULL);
    if (!ctx) {
        fprintf(stderr, "ERROR: OSMesaCreateContext failed");
    } else if (!OSMesaMakeCurrent(ctx, result.pixel, GL_UNSIGNED_BYTE, result.width, result.height)) {
        fprintf(stderr, "ERROR: OSMesaMakeCurrent failed");
        OSMesaDestroyContext(ctx);
        ctx = NULL;
    }

    return ctx;
}

void bind(Texture texture) {
    glBindTexture(GL_TEXTURE_2D, texture);
}

void upload(Texture texture, Image image) {
    bind(texture);
    glTexImage2D(GL_TEXTURE_2D, 0, GL_RGBA8, image.width, image.height, 0, GL_RGBA, GL_UNSIGNED_BYTE, image.pixel);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST);
}

void draw(float x, float y, float z, TextureType type) {
    glPushMatrix();
    glBegin(GL_QUADS);
    glNormal3f(0.0f, 0.0f, -1.0f);

    if (type == TEXTURE_TYPE_NONE) {
        int i;
        for (i = 0; i < vertex_count; i += 3) {
            glVertex3f(vertices[i] * x, vertices[i + 1] * y, vertices[i + 2] * z);
        }
    } else {
        int i;
        for (i = 0; i < 24; i++) {
            glTexCoord2f(u[type][i], v[type][i]);
            int idx = i * 3;
            glVertex3f(vertices[idx] * x, vertices[idx + 1] * y, vertices[idx + 2] * z);
        }
    }

    glEnd();
    glPopMatrix();
}

void renderPlayer(bool overlay, bool body) {
    float offset = overlay ? 0.05f : 0.0f;

    if (body) {
        glPushMatrix();
        glTranslatef(0.0f, -2.5f, 0.0f);
        bind(overlay ? TEXTURE_BODY_OVERLAY : TEXTURE_BODY);
        draw(1.0f + offset, 1.5f + offset, 0.5f + offset, TEXTURE_TYPE_TORSO);


        glPushMatrix();
        bind(overlay ? TEXTURE_LEFT_LEG_OVERLAY : TEXTURE_LEFT_LEG);
        glTranslatef(0.5f, -3.0f, 0.0f);
        draw(0.5f + offset, 1.5f + offset, 0.5f + offset, TEXTURE_TYPE_LIMB);
        glPopMatrix();

        glPushMatrix();
        bind(overlay ? TEXTURE_RIGHT_LEG_OVERLAY : TEXTURE_RIGHT_LEG);
        glTranslatef(-0.5f, -3.0f, 0.0f);
        draw(0.5f + offset, 1.5f + offset, 0.5f + offset, TEXTURE_TYPE_LIMB);
        glPopMatrix();

        glPushMatrix();
        bind(overlay ? TEXTURE_LEFT_ARM_OVERLAY : TEXTURE_LEFT_ARM);
        glTranslatef(1.75f, 0.1f, 0.0f);
        glRotatef(10.0f, 0.0f, 0.0f, 1.0f);
        draw(0.5f + offset, 1.5f + offset, 0.5f + offset, TEXTURE_TYPE_LIMB);
        glPopMatrix();

        glPushMatrix();
        bind(overlay ? TEXTURE_RIGHT_ARM_OVERLAY : TEXTURE_RIGHT_ARM);
        glTranslatef(-1.75f, 0.1f, 0.0f);
        glRotatef(-10.0f, 0.0f, 0.0f, 1.0f);
        draw(0.5f + offset, 1.5f + offset, 0.5f + offset, TEXTURE_TYPE_LIMB);
        glPopMatrix();

        glPopMatrix();
    }

    bind(overlay ? TEXTURE_HEAD_OVERLAY : TEXTURE_HEAD);
    draw(1.0f + offset, 1.0f + offset, 1.0f + offset, TEXTURE_TYPE_HEAD);
}

bool Render(
        float angle, float tilt, float zoom,
        bool shadow, bool lighting,
        bool portrait, bool full,
        bool overlay, bool newSkin,
        Image result,
        Image head, Image *headOverlay,
        Image *body, Image *bodyOverlay,
        Image *leftArm, Image *leftArmOverlay, Image *rightArm, Image *rightArmOverlay,
        Image *leftLeg, Image *leftLegOverlay, Image *rightLeg, Image *rightLegOverlay) {

    if (full) portrait = true;

    OSMesaContext ctx = initOSMesa(result);

    if (ctx) {
        // Initalize OpenGL
        glClearColor(0.0f, 0.0f, 0.0f, 0.0f);
        glClearDepth(1.0f);
        glShadeModel(GL_SMOOTH);
        glEnable(GL_DEPTH_TEST);
        glDepthFunc(GL_LEQUAL);

        glMatrixMode(GL_PROJECTION);
        glLoadIdentity();

        gluPerspective(45.0f, (float) result.width / (float) result.height, 0.1f, 100.0f);
        glHint(GL_PERSPECTIVE_CORRECTION_HINT, GL_NICEST);

        glMatrixMode(GL_MODELVIEW);
        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);

        glPushMatrix();
        glCullFace(GL_BACK);

        // Upload head
        upload(TEXTURE_HEAD, head);
        if (overlay)
            upload(TEXTURE_HEAD_OVERLAY, *headOverlay);

        if (portrait) {
            upload(TEXTURE_BODY, *body);
            upload(TEXTURE_LEFT_ARM, *leftArm);
            upload(TEXTURE_RIGHT_ARM, *rightArm);
            upload(TEXTURE_LEFT_LEG, *leftLeg);
            upload(TEXTURE_RIGHT_LEG, *rightLeg);

            if (overlay) {
                upload(TEXTURE_BODY_OVERLAY, *bodyOverlay);
                upload(TEXTURE_LEFT_ARM_OVERLAY, *leftArmOverlay);
                upload(TEXTURE_RIGHT_ARM_OVERLAY, *rightArmOverlay);
                upload(TEXTURE_LEFT_LEG_OVERLAY, *leftLegOverlay);
                upload(TEXTURE_RIGHT_LEG_OVERLAY, *rightLegOverlay);
            }
        }

        if (full)
            glTranslatef(0.0f, 2.5f, -4.5f);
        else if (portrait)
            glTranslatef(0.0f, 1.0f, -2.0f);
        else
            glTranslatef(0.0f, 0.20f, 0.0f);

        glTranslatef(0, 0, zoom);
        glRotatef(tilt, 1.0f, 0.0f, 0.0f);
        glRotatef(angle, 0.0f, 1.0f, 0.0f);

        glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA);
        glEnable(GL_BLEND);
        glEnable(GL_ALPHA_TEST);
        glAlphaFunc(GL_SRC_ALPHA, GL_GREATER);
        glEnable(GL_DEPTH_TEST);

        if (shadow) {
            // Setup shadow
            glPushMatrix();

            glTranslatef(0.0f, -1.0f, 0.0f);
            GLfloat scaleX = 0.97f;
            GLfloat scaleZ = 0.97f;
            if (portrait) {
                scaleZ = 0.47f;
                glTranslatef(0.0f, -6.05f, 0.0f);
            }

            static const GLfloat count = 10;
            static const GLfloat inc = 0.02f;

            GLfloat i;
            for (i = 0; i < count; i++) {
                scaleX += inc;
                scaleZ += inc;

                glTranslatef(0.0f, -0.001f, 0.0f);
                glColor4f(0.0f, 0.0f, 0.0f, (1.0f - (i / (float) count)) / 2.0f);
                draw(scaleX, 0.01f, scaleZ, TEXTURE_TYPE_NONE);
            }

            glPopMatrix();
        }

        if (lighting) {
            // Setup lighting
            glEnable(GL_LIGHTING);
            glEnable(GL_LIGHT0);

            static const float position[4] = {-4.0f, 2.0f, 1.0f, 100.0f};
            static const float ambient[4] = {3.0f, 3.0f, 3.0f, 1.0f};
            glLightfv(GL_LIGHT0, GL_POSITION, position);
            glLightfv(GL_LIGHT0, GL_AMBIENT, ambient);
        }

        glEnable(GL_TEXTURE_2D);
        glColor3f(1, 1, 1);
        renderPlayer(false, portrait);
        if (overlay) {
            renderPlayer(true, portrait && newSkin);
        }

        glPopMatrix();
        glFinish();
        OSMesaDestroyContext(ctx);
        return true;
    }

    return false;
}
