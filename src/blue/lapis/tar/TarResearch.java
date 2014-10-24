package blue.lapis.tar;

import java.awt.image.BufferedImage;
import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.InputStream;
import java.net.URL;
import java.nio.ByteBuffer;

import javax.imageio.ImageIO;

import org.lwjgl.BufferUtils;
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
			upload(head);
			draw();
			BufferedImage img = readPixels();
			ImageIO.write(img, "png", new File(player+".png"));
		} catch (Exception e) {
			e.printStackTrace();
			System.exit(1);
		}
	}
	private static BufferedImage readPixels() {
		GL11.glReadBuffer(GL11.GL_FRONT);
		ByteBuffer buf = BufferUtils.createByteBuffer(width * height * 4);
		GL11.glReadPixels(0, 0, width, height, GL11.GL_RGBA, GL11.GL_UNSIGNED_BYTE, buf);
		BufferedImage img = new BufferedImage(width, height, BufferedImage.TYPE_INT_ARGB);
		  
		for(int x = 0; x < width; x++)  {
			for(int y = 0; y < height; y++) {
				int i = (x + (width * y)) * 4;
				int r = buf.get(i) & 0xFF;
				int g = buf.get(i + 1) & 0xFF;
				int b = buf.get(i + 2) & 0xFF;
				int a = buf.get(i + 3) & 0xFF;
				img.setRGB(x, height - (y + 1), (a << 24) | (r << 16) | (g << 8) | b);
			}
		}
		return img;
	}
	private static float quadRot = 45;
	private static void draw() {
		GL11.glClear(GL11.GL_COLOR_BUFFER_BIT | GL11.GL_DEPTH_BUFFER_BIT);

        GL11.glLoadIdentity();
        
        GL11.glLoadIdentity();
        GL11.glTranslatef(0,0,-5);
        GL11.glRotatef(quadRot,0f,1.0f,0f);
        GL11.glColor3f(1, 1, 1);
        GL11.glBegin(GL11.GL_QUADS);
        	GL11.glTexCoord2f(0f, 0f);
            GL11.glVertex3f( 1.0f, 1.0f,-1.0f);
            GL11.glVertex3f(-1.0f, 1.0f,-1.0f);
            GL11.glVertex3f(-1.0f, 1.0f, 1.0f);
            GL11.glVertex3f( 1.0f, 1.0f, 1.0f);
            GL11.glTexCoord2f(1f, 0f);
            GL11.glVertex3f( 1.0f,-1.0f, 1.0f);
            GL11.glVertex3f(-1.0f,-1.0f, 1.0f);
            GL11.glVertex3f(-1.0f,-1.0f,-1.0f);
            GL11.glVertex3f( 1.0f,-1.0f,-1.0f);
            GL11.glTexCoord2f(1f, 1f);
            GL11.glVertex3f( 1.0f, 1.0f, 1.0f);
            GL11.glVertex3f(-1.0f, 1.0f, 1.0f);
            GL11.glVertex3f(-1.0f,-1.0f, 1.0f);
            GL11.glVertex3f( 1.0f,-1.0f, 1.0f);

            GL11.glVertex3f( 1.0f,-1.0f,-1.0f);
            GL11.glVertex3f(-1.0f,-1.0f,-1.0f);
            GL11.glVertex3f(-1.0f, 1.0f,-1.0f);
            GL11.glVertex3f( 1.0f, 1.0f,-1.0f);

            GL11.glVertex3f(-1.0f, 1.0f, 1.0f);
            GL11.glVertex3f(-1.0f, 1.0f,-1.0f);
            GL11.glVertex3f(-1.0f,-1.0f,-1.0f);
            GL11.glVertex3f(-1.0f,-1.0f, 1.0f);

            GL11.glVertex3f( 1.0f, 1.0f,-1.0f);
            GL11.glVertex3f( 1.0f, 1.0f, 1.0f);
            GL11.glVertex3f( 1.0f,-1.0f, 1.0f);
            GL11.glVertex3f( 1.0f,-1.0f,-1.0f);
        GL11.glEnd();
	}
	private static void upload(BufferedImage img) {
		int[] pixels = new int[img.getWidth()*img.getHeight()];
		img.getRGB(0, 0, img.getWidth(), img.getHeight(), pixels, 0, img.getWidth());
		ByteBuffer buf = BufferUtils.createByteBuffer(img.getWidth() * img.getHeight() * 4);
		for (int y = 0; y < img.getHeight(); y++) {
			for (int x = 0; x < img.getWidth(); x++) {
				int pixel = pixels[y*img.getWidth()+x];
				buf.put((byte) ((pixel >> 16) & 0xFF));
				buf.put((byte) ((pixel >> 8) & 0xFF));
                buf.put((byte) (pixel & 0xFF));
                buf.put((byte) ((pixel >> 24) & 0xFF));
			}
		}
		buf.flip();
		
		GL11.glTexImage2D(GL11.GL_TEXTURE_2D, 0, GL11.GL_RGBA8, img.getWidth(), img.getHeight(), 0, GL11.GL_RGBA, GL11.GL_UNSIGNED_BYTE, buf);
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
