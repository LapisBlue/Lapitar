#include <stdio.h>
#include <stdbool.h>
#include <GL/osmesa.h>
#include <GL/glu.h>
#include "head.h"

OSMesaContext init_osmesa(int width, int height, void *buffer) {
    OSMesaContext ctx = OSMesaCreateContextExt(OSMESA_RGBA, 16, 0, 0, NULL);
    if (!ctx) {
        fprintf(stderr, "ERROR: OSMesaCreateContext failed");
    } else if (!OSMesaMakeCurrent(ctx, buffer, GL_UNSIGNED_BYTE, width, height)) {
        fprintf(stderr, "ERROR: OSMesaMakeCurrent failed");
        OSMesaDestroyContext(ctx);
        ctx = NULL;
    }

    return ctx;
}

void init_gl(int width, int height) {
    glClearColor(0.0f, 0.0f, 0.0f, 0.0f);
    glClearDepth(0.0f);
    glShadeModel(GL_SMOOTH);
    glEnable(GL_DEPTH_TEST);
    glDepthFunc(GL_LEQUAL);

    glMatrixMode(GL_PROJECTION);
    glLoadIdentity();

    gluPerspective(45.0f, (float) width / (float) height, 0.1f, 100.0f);
    glMatrixMode(GL_MODELVIEW);

    glHint(GL_PERSPECTIVE_CORRECTION_HINT, GL_NICEST);
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
}

void draw(float x, float y, float z, float angle) {
    glPushMatrix();
    glRotatef(20.0f, 1.0f, 0.0f, 0.0f);
    glTranslatef(0.0f, -1.5f, -4.5f);
    glRotatef(angle, 0.0f, 1.0f, 0.0f);

    // Lighting
    if (glIsEnabled(GL_LIGHTING)) {
        const float position[4] = {-4.0f, 2.0f, 1.0f, 100.0f};
        const float ambient[4] = {3.0f, 3.0f, 3.0f, 1.0f};
        glLightfv(GL_LIGHT0, GL_POSITION, position);
        glLightfv(GL_LIGHT0, GL_AMBIENT, ambient);
    }

    glBegin(GL_QUADS);
    glNormal3f(0.0f, 0.0f, -1.0f);

    // Front
    glTexCoord2f(0.25f, 1.0f);
    glVertex3f(-x, -y, z);
    glTexCoord2f(0.5f, 1.0f);
    glVertex3f(x, -y, z);
    glTexCoord2f(0.5f, 0.5f);
    glVertex3f(x, y, z);
    glTexCoord2f(0.25f, 0.5f);
    glVertex3f(-x, y, z);

    // Back
    glTexCoord2f(1.0f, 1.0f);
    glVertex3f(-x, -y, -z);
    glTexCoord2f(1.0f, 0.5f);
    glVertex3f(-x, y, -z);
    glTexCoord2f(0.75f, 0.5f);
    glVertex3f(x, y, -z);
    glTexCoord2f(0.75f, 1.0f);
    glVertex3f(x, -y, -z);

    // Top
    glTexCoord2f(0.5f, 0.0f);
    glVertex3f(-x, y, -z);
    glTexCoord2f(0.5f, 0.5f);
    glVertex3f(-x, y, z);
    glTexCoord2f(0.25f, 0.5f);
    glVertex3f(x, y, z);
    glTexCoord2f(0.25f, 0.0f);
    glVertex3f(x, y, -z);

    // Bottom
    glTexCoord2f(0.5f, 0.5f);
    glVertex3f(-x, -y, -z);
    glTexCoord2f(0.75f, 0.5f);
    glVertex3f(x, -y, -z);
    glTexCoord2f(0.75f, 0.0f);
    glVertex3f(x, -y, z);
    glTexCoord2f(0.5f, 0.0f);
    glVertex3f(-x, -y, z);

    // Left
    glTexCoord2f(0.75f, 1.0f);
    glVertex3f(x, -y, -z);
    glTexCoord2f(0.75f, 0.5f);
    glVertex3f(x, y, -z);
    glTexCoord2f(0.5f, 0.5f);
    glVertex3f(x, y, z);
    glTexCoord2f(0.5f, 1.0f);
    glVertex3f(x, -y, z);

    // Right
    glTexCoord2f(0.0f, 1.0f);
    glVertex3f(-x, -y, -z);
    glTexCoord2f(0.25f, 1.0f);
    glVertex3f(-x, -y, z);
    glTexCoord2f(0.25f, 0.5f);
    glVertex3f(-x, y, z);
    glTexCoord2f(0.0f, 0.5f);
    glVertex3f(-x, y, -z);

    glEnd();
    glPopMatrix();
}

void setup_shadow(float angle) {
    glEnable(GL_BLEND);
    glDisable(GL_TEXTURE_2D);
    glPushMatrix();
    glTranslatef(0.0f, -0.95f, -0.45f);

    const GLfloat count = 10.0f;
    GLfloat i;
    for (i = 0; i < count; i++) {
        glTranslatef(0.0f, -0.01f, 0.0f);
        glColor4f(0.0f, 0.0f, 0.0f, (1.0f - (i / count)) / 2.0f);
        draw(1.02f, 0.01f, 1.02f, angle);
    }

    glPopMatrix();
}

void setup_lighting(void) {
    glEnable(GL_LIGHTING);
    glEnable(GL_LIGHT0);
}

void upload_image(unsigned int id, void *buffer, int width, int height) {
    glBindTexture(GL_TEXTURE_2D, id);
    glTexImage2D(GL_TEXTURE_2D, 0, GL_RGBA8, width, height, 0, GL_RGBA, GL_UNSIGNED_BYTE, buffer);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST);
}

bool RenderHead(
        float angle,
        bool shadow, bool lighting,
        void *result, int width, int height,
        void *head, int headWidth, int headHeight,
        void *helm, int helmWidth, int helmHeight) {
    OSMesaContext ctx = init_osmesa(width, height, result);

    if (ctx) {
        // Initialize OpenGL
        init_gl(width, height);

        // Upload helm
        upload_image(1, head, headWidth, headHeight);
        if (helm)
            upload_image(2, helm, helmWidth, helmHeight);

        if (shadow)
            setup_shadow(angle);

        glEnable(GL_TEXTURE_2D);
        if (lighting)
            setup_lighting();

        glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA);
        glColor3f(1.0f, 1.0f, 1.0f);

        glBindTexture(GL_TEXTURE_2D, 1);
        draw(1.0f, 1.0f, 1.0f, angle);

        if (helm) {
            glBindTexture(GL_TEXTURE_2D, 2);
            draw(1.05f, 1.05f, 1.05f, angle);
        }

        OSMesaDestroyContext(ctx);
        return true;
    }

    return false;
}
