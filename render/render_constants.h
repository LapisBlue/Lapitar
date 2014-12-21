typedef enum Texture {
    TEXTURE_HEAD = 1, TEXTURE_HEAD_OVERLAY,
    TEXTURE_BODY, TEXTURE_BODY_OVERLAY,
    TEXTURE_LEFT_ARM, TEXTURE_LEFT_ARM_OVERLAY,
    TEXTURE_RIGHT_ARM, TEXTURE_RIGHT_ARM_OVERLAY,
    TEXTURE_LEFT_LEG, TEXTURE_LEFT_LEG_OVERLAY,
    TEXTURE_RIGHT_LEG, TEXTURE_RIGHT_LEG_OVERLAY
} Texture;

typedef enum TextureType {
    TEXTURE_TYPE_NONE, TEXTURE_TYPE_HEAD, TEXTURE_TYPE_TORSO, TEXTURE_TYPE_LIMB
} TextureType;

static const float u[3][24] = {
        { // Head
                8.0f / 32.0f, 16.0f / 32.0f, 16.0f / 32.0f, 8.0f / 32.0f,
                24.0f / 32.0f, 32.0f / 32.0f, 32.0f / 32.0f, 24.0f / 32.0f,
                8.0f / 32.0f, 16.0f / 32.0f, 16.0f / 32.0f, 8.0f / 32.0f,
                16.0f / 32.0f, 24.0f / 32.0f, 24.0f / 32.0f, 16.0f / 32.0f,
                16.0f / 32.0f, 24.0f / 32.0f, 24.0f / 32.0f, 16.0f / 32.0f,
                0.0f / 32.0f, 8.0f / 32.0f, 8.0f / 32.0f, 0.0f / 32.0f
        },
        { // Torso
                4.0f / 24.0f, 12.0f / 24.0f, 12.0f / 24.0f, 4.0f / 24.0f,
                16.0f / 24.0f, 24.0f / 24.0f, 24.0f / 24.0f, 16.0f / 24.0f,
                4.0f / 24.0f, 12.0f / 24.0f, 12.0f / 24.0f, 4.0f / 24.0f,
                12.0f / 24.0f, 20.0f / 24.0f, 20.0f / 24.0f, 12.0f / 24.0f,
                12.0f / 24.0f, 16.0f / 24.0f, 16.0f / 24.0f, 12.0f / 24.0f,
                0.0f / 24.0f, 4.0f / 24.0f, 4.0f / 24.0f, 0.0f / 24.0f
        },
        { // Limb
                4.0f / 16.0f, 8.0f / 16.0f, 8.0f / 16.0f, 4.0f / 16.0f,
                12.0f / 16.0f, 16.0f / 16.0f, 16.0f / 16.0f, 12.0f / 16.0f,
                4.0f / 16.0f, 8.0f / 16.0f, 8.0f / 16.0f, 4.0f / 16.0f,
                8.0f / 16.0f, 12.0f / 16.0f, 12.0f / 16.0f, 8.0f / 16.0f,
                8.0f / 16.0f, 12.0f / 16.0f, 12.0f / 16.0f, 8.0f / 16.0f,
                0.0f / 16.0f, 4.0f / 16.0f, 4.0f / 16.0f, 0.0f / 16.0f
        }
};

static const float v[3][24] = {
        { // Head
                16.0f / 16.0f, 16.0f / 16.0f, 8.0f / 16.0f, 8.0f / 16.0f,
                16.0f / 16.0f, 16.0f / 16.0f, 8.0f / 16.0f, 8.0f / 16.0f,
                8.0f / 16.0f, 8.0f / 16.0f, 0.0f / 16.0f, 0.0f / 16.0f,
                8.0f / 16.0f, 8.0f / 16.0f, 0.0f / 16.0f, 0.0f / 16.0f,
                16.0f / 16.0f, 16.0f / 16.0f, 8.0f / 16.0f, 8.0f / 16.0f,
                16.0f / 16.0f, 16.0f / 16.0f, 8.0f / 16.0f, 8.0f / 16.0f
        },
        { // Torso
                16.0f / 16.0f, 16.0f / 16.0f, 4.0f / 16.0f, 4.0f / 16.0f,
                16.0f / 16.0f, 16.0f / 16.0f, 4.0f / 16.0f, 4.0f / 16.0f,
                4.0f / 16.0f, 4.0f / 16.0f, 0.0f / 16.0f, 0.0f / 16.0f,
                4.0f / 16.0f, 4.0f / 16.0f, 0.0f / 16.0f, 0.0f / 16.0f,
                16.0f / 16.0f, 16.0f / 16.0f, 4.0f / 16.0f, 4.0f / 16.0f,
                16.0f / 16.0f, 16.0f / 16.0f, 4.0f / 16.0f, 4.0f / 16.0f
        },
        { // Limb
                16.0f / 16.0f, 16.0f / 16.0f, 4.0f / 16.0f, 4.0f / 16.0f,
                16.0f / 16.0f, 16.0f / 16.0f, 4.0f / 16.0f, 4.0f / 16.0f,
                4.0f / 16.0f, 4.0f / 16.0f, 0.0f / 16.0f, 0.0f / 16.0f,
                4.0f / 16.0f, 4.0f / 16.0f, 0.0f / 16.0f, 0.0f / 16.0f,
                16.0f / 16.0f, 16.0f / 16.0f, 4.0f / 16.0f, 4.0f / 16.0f,
                16.0f / 16.0f, 16.0f / 16.0f, 4.0f / 16.0f, 4.0f / 16.0f
        }
};

static const int vertice_count = 72;

static const float vertices[72] = {

        // Front
        -1.0f, -1.0f,  1.0f,
        1.0f, -1.0f,  1.0f,
        1.0f,  1.0f,  1.0f,
        -1.0f,  1.0f,  1.0f,

        // Back
        -1.0f, -1.0f, -1.0f,
        1.0f, -1.0f, -1.0f,
        1.0f,  1.0f, -1.0f,
        -1.0f,  1.0f, -1.0f,

        // Top
        -1.0f,  1.0f,  1.0f,
        1.0f,  1.0f,  1.0f,
        1.0f,  1.0f, -1.0f,
        -1.0f,  1.0f, -1.0f,

        // Bottom
        -1.0f, -1.0f, -1.0f,
        1.0f, -1.0f, -1.0f,
        1.0f, -1.0f,  1.0f,
        -1.0f, -1.0f,  1.0f,

        // Left
        1.0f, -1.0f,  1.0f,
        1.0f, -1.0f, -1.0f,
        1.0f,  1.0f, -1.0f,
        1.0f,  1.0f,  1.0f,

        // Right
        -1.0f, -1.0f, -1.0f,
        -1.0f, -1.0f,  1.0f,
        -1.0f,  1.0f,  1.0f,
        -1.0f,  1.0f, -1.0f
};
