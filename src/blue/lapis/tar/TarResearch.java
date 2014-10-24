package blue.lapis.tar;

import java.awt.image.BufferedImage;
import java.io.ByteArrayOutputStream;
import java.io.InputStream;
import java.net.URL;

import javax.imageio.ImageIO;

import org.lwjgl.opengl.Display;
import org.lwjgl.opengl.DisplayMode;
import org.lwjgl.opengl.GL11;
import org.lwjgl.opengl.Pbuffer;
import org.lwjgl.opengl.PixelFormat;
import org.lwjgl.util.glu.GLU;

public class TarResearch {
	// make these switches at some point
	public static final int width = 256;
	public static final int height = 256;
	private static Pbuffer buffer;
	
	public static void main(String[] args) {
		String player;
		if (args.length > 0) {
			player = args[0];
		} else {
			player = "Aesen_";
		}
		System.out.println("Creating avatar for "+player);
		try {
			init();
			BufferedImage skin = ImageIO.read(new URL("https://s3.amazonaws.com/MinecraftSkins/"+player+".png"));
			BufferedImage head = skin.getSubimage(0, 0, 32, 16);
			upload(skin, 1);
			draw(1);
		} catch (Exception e) {
			e.printStackTrace();
			System.exit(1);
		}
	}
	private static void upload(BufferedImage skin, int glId) {
		
	}
	private static void createBuffer() throws Exception {
		buffer = new Pbuffer(width, height, new PixelFormat(), null, null);
		buffer.makeCurrent();
		if (buffer.isBufferLost()) {
			cleanup();
			System.err.println("Failed to set up PBuffer.");
			System.exit(2);
		}
    }
    private static void init() throws Exception {
        createBuffer();
        initGL();
    }

    private static void initGL() {
        GL11.glEnable(GL11.GL_TEXTURE_2D);
        GL11.glShadeModel(GL11.GL_SMOOTH);
        GL11.glClearColor(0.0f, 0.0f, 0.0f, 0.0f);
        GL11.glClearDepth(1.0);
        GL11.glEnable(GL11.GL_DEPTH_TEST);
        GL11.glDepthFunc(GL11.GL_LEQUAL);

        GL11.glMatrixMode(GL11.GL_PROJECTION);
        GL11.glLoadIdentity();

        GLU.gluPerspective(
          45.0f,
          (float)width / (float)height,
          0.1f,
          100.0f);
        GL11.glMatrixMode(GL11.GL_MODELVIEW);

        GL11.glHint(GL11.GL_PERSPECTIVE_CORRECTION_HINT, GL11.GL_NICEST);
    }
    private static void cleanup() {
        buffer.destroy();
    }
}
