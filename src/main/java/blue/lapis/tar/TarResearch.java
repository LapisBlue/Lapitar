package blue.lapis.tar;

import java.awt.Graphics2D;
import java.awt.Image;
import java.awt.image.BufferedImage;
import java.io.File;
import java.net.URL;
import java.nio.ByteBuffer;
import java.nio.FloatBuffer;
import java.nio.IntBuffer;

import javax.imageio.ImageIO;

import org.lwjgl.BufferUtils;
import org.lwjgl.opengl.GL11;
import org.lwjgl.opengl.Pbuffer;
import org.lwjgl.opengl.PixelFormat;
import org.lwjgl.util.glu.GLU;

public class TarResearch {
	// make these switches at some point
	public static final int supersampling = 4;
	public static final int width = 256*supersampling;
	public static final int height = 256*supersampling;
	private static Pbuffer buffer;
	private static FloatBuffer lightPosition;
	private static FloatBuffer lightAmbient;
	
	public static void main(String[] args) {
		if (args.length == 0) {
			args = new String[] { "Aesen_" };
		}
		for (String player : args) {
			System.out.println("Creating avatar for "+player);
			try {
				init();
				BufferedImage skin = ImageIO.read(new URL("https://s3.amazonaws.com/MinecraftSkins/"+player+".png"));
				BufferedImage head = skin.getSubimage(0, 0, 32, 16);
				ImageIO.write(head, "png", new File(player+".head.png"));
				BufferedImage helm = skin.getSubimage(32, 0, 32, 16);
				ImageIO.write(helm, "png", new File(player+".helm.png"));
				GL11.glClear(GL11.GL_COLOR_BUFFER_BIT | GL11.GL_DEPTH_BUFFER_BIT);
				lightPosition = BufferUtils.createFloatBuffer(4);
			    lightPosition.mark();
			    lightPosition.put(-4f);
			    lightPosition.put(2f);
			    lightPosition.put(1f);
			    lightPosition.put(100f);
			    lightPosition.reset();
			    
			    lightAmbient = BufferUtils.createFloatBuffer(4);
			    lightAmbient.mark();
			    lightAmbient.put(3.0f);
			    lightAmbient.put(3.0f);
			    lightAmbient.put(3.0f);
			    lightAmbient.put(1f);
			    lightAmbient.reset();
				upload(head, 1);
				upload(helm, 2);
				
				GL11.glEnable(GL11.GL_BLEND);
				GL11.glDisable(GL11.GL_TEXTURE_2D);
				GL11.glPushMatrix();
					GL11.glTranslatef(0f, -0.95f, -0.45f);
					int count = 10;
					for (int i = 0; i < count; i++) {
						GL11.glTranslatef(0f, -0.01f, 0f);
						GL11.glColor4f(0, 0, 0, (1-(i/(float)count))/2f);
						draw(1.02f, 0.01f, 1.02f);
					}
				GL11.glPopMatrix();
	
				GL11.glEnable(GL11.GL_TEXTURE_2D);
				GL11.glEnable(GL11.GL_LIGHTING);
			    GL11.glEnable(GL11.GL_LIGHT0);
	
			    GL11.glBlendFunc(GL11.GL_SRC_ALPHA, GL11.GL_ONE_MINUS_SRC_ALPHA);
				GL11.glColor3f(1, 1, 1);
				GL11.glBindTexture(GL11.GL_TEXTURE_2D, 1);
					draw(1.0f, 1.0f, 1.0f);
				GL11.glBindTexture(GL11.GL_TEXTURE_2D, 2);
					draw(1.05f, 1.05f, 1.05f);
				
				BufferedImage img = readPixels();
				BufferedImage out = new BufferedImage(width/supersampling, height/supersampling, BufferedImage.TYPE_INT_ARGB);
				Graphics2D gout = out.createGraphics();
				gout.drawImage(img.getScaledInstance(img.getWidth()/supersampling, img.getHeight()/supersampling, Image.SCALE_SMOOTH), 0, 0, null);
				gout.dispose();
				ImageIO.write(out, "png", new File(player+".png"));
				cleanup();
			} catch (Exception e) {
				e.printStackTrace();
				System.exit(1);
			}
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
	private static float quadRot = 0;
	private static void draw(float xScale, float yScale, float zScale) throws Exception {
		GL11.glPushMatrix();
			GL11.glTranslatef(0,0.25f,-5);
			GL11.glRotatef(quadRot,0f,1.0f,0f);
			GL11.glRotatef(20,1.0f,0f,0.0f);
			GL11.glLight(GL11.GL_LIGHT0, GL11.GL_POSITION, lightPosition);
			GL11.glLight(GL11.GL_LIGHT0, GL11.GL_AMBIENT, lightAmbient);
			GL11.glBegin(GL11.GL_QUADS);
				GL11.glNormal3f(0, 0, -1f);
				// Front
				GL11.glTexCoord2f( 0.25f, 1.00f ); GL11.glVertex3f( -1.0f*xScale, -1.0f*yScale, 1.0f*zScale);
				GL11.glTexCoord2f( 0.50f, 1.00f ); GL11.glVertex3f(  1.0f*xScale, -1.0f*yScale, 1.0f*zScale);
				GL11.glTexCoord2f( 0.50f, 0.50f ); GL11.glVertex3f(  1.0f*xScale,  1.0f*yScale, 1.0f*zScale);
				GL11.glTexCoord2f( 0.25f, 0.50f ); GL11.glVertex3f( -1.0f*xScale,  1.0f*yScale, 1.0f*zScale);
				// Back
				GL11.glTexCoord2f( 1.00f, 1.00f ); GL11.glVertex3f( -1.0f*xScale, -1.0f*yScale, -1.0f*zScale);
				GL11.glTexCoord2f( 1.00f, 0.50f ); GL11.glVertex3f( -1.0f*xScale,  1.0f*yScale, -1.0f*zScale);
				GL11.glTexCoord2f( 0.75f, 0.50f ); GL11.glVertex3f(  1.0f*xScale,  1.0f*yScale, -1.0f*zScale);
				GL11.glTexCoord2f( 0.75f, 1.00f ); GL11.glVertex3f(  1.0f*xScale, -1.0f*yScale, -1.0f*zScale);
				// Top
				GL11.glTexCoord2f( 0.50f, 0.00f ); GL11.glVertex3f( -1.0f*xScale,  1.0f*yScale, -1.0f*zScale);
				GL11.glTexCoord2f( 0.50f, 0.50f ); GL11.glVertex3f( -1.0f*xScale,  1.0f*yScale,  1.0f*zScale);
				GL11.glTexCoord2f( 0.25f, 0.50f ); GL11.glVertex3f(  1.0f*xScale,  1.0f*yScale,  1.0f*zScale);
				GL11.glTexCoord2f( 0.25f, 0.00f ); GL11.glVertex3f(  1.0f*xScale,  1.0f*yScale, -1.0f*zScale);
				// Bottom
				GL11.glTexCoord2f( 0.50f, 0.50f ); GL11.glVertex3f( -1.0f*xScale, -1.0f*yScale, -1.0f*zScale);
				GL11.glTexCoord2f( 0.75f, 0.50f ); GL11.glVertex3f(  1.0f*xScale, -1.0f*yScale, -1.0f*zScale);
				GL11.glTexCoord2f( 0.75f, 0.00f ); GL11.glVertex3f(  1.0f*xScale, -1.0f*yScale,  1.0f*zScale);
				GL11.glTexCoord2f( 0.50f, 0.00f ); GL11.glVertex3f( -1.0f*xScale, -1.0f*yScale,  1.0f*zScale);
				// Left
				GL11.glTexCoord2f( 0.75f, 1.00f ); GL11.glVertex3f(  1.0f*xScale, -1.0f*yScale, -1.0f*zScale);
				GL11.glTexCoord2f( 0.75f, 0.50f ); GL11.glVertex3f(  1.0f*xScale,  1.0f*yScale, -1.0f*zScale);
				GL11.glTexCoord2f( 0.50f, 0.50f ); GL11.glVertex3f(  1.0f*xScale,  1.0f*yScale,  1.0f*zScale);
				GL11.glTexCoord2f( 0.50f, 1.00f ); GL11.glVertex3f(  1.0f*xScale, -1.0f*yScale,  1.0f*zScale);
				// Right
				GL11.glTexCoord2f( 0.00f, 1.00f ); GL11.glVertex3f( -1.0f*xScale, -1.0f*yScale, -1.0f*zScale);
				GL11.glTexCoord2f( 0.25f, 1.00f ); GL11.glVertex3f( -1.0f*xScale, -1.0f*yScale,  1.0f*zScale);
				GL11.glTexCoord2f( 0.25f, 0.50f ); GL11.glVertex3f( -1.0f*xScale,  1.0f*yScale,  1.0f*zScale);
				GL11.glTexCoord2f( 0.00f, 0.50f ); GL11.glVertex3f( -1.0f*xScale,  1.0f*yScale, -1.0f*zScale);
			GL11.glEnd();
		GL11.glPopMatrix();
	}
	private static int upload(BufferedImage img, int id) {
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
		GL11.glBindTexture(GL11.GL_TEXTURE_2D, id);
		GL11.glTexImage2D(GL11.GL_TEXTURE_2D, 0, GL11.GL_RGBA8, img.getWidth(), img.getHeight(), 0, GL11.GL_RGBA, GL11.GL_UNSIGNED_BYTE, buf);
		GL11.glTexParameteri(GL11.GL_TEXTURE_2D, GL11.GL_TEXTURE_MAG_FILTER, GL11.GL_NEAREST);
		GL11.glTexParameteri(GL11.GL_TEXTURE_2D, GL11.GL_TEXTURE_MIN_FILTER, GL11.GL_NEAREST);
		return id;
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
		GL11.glClearColor(0.0f, 0.0f, 0.0f, 0.0f);
		GL11.glClearDepth(1.0);
		GL11.glShadeModel(GL11.GL_SMOOTH);
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
