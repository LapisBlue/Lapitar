package blue.lapis.tar;

import java.awt.Graphics2D;
import java.awt.Image;
import java.awt.image.BufferedImage;
import java.nio.ByteBuffer;
import java.nio.FloatBuffer;

import org.lwjgl.BufferUtils;
import org.lwjgl.LWJGLException;
import org.lwjgl.opengl.Display;
import org.lwjgl.opengl.DisplayMode;
import org.lwjgl.opengl.GL11;
import org.lwjgl.opengl.Pbuffer;
import org.lwjgl.opengl.PixelFormat;
import org.lwjgl.util.glu.GLU;

import com.google.common.base.MoreObjects;

public class TarRenderer {
	private float angle;
	private final int width, height, finalWidth, finalHeight;
	private final int superSampling;

	private final boolean helmet;
	private final boolean shadow;
	private final boolean lighting;

	private final boolean useWindow;
	
	private Pbuffer buffer;
	private FloatBuffer lightPosition;
	private FloatBuffer lightAmbient;
	
	private boolean initialized = false;

	public TarRenderer(float angle, int width, int height, int superSampling, boolean helmet, boolean shadow, boolean lighting, boolean useWindow) {
		this.angle = angle;
		this.width = width * superSampling;
		this.finalWidth = width;
		this.height = height * superSampling;
		this.finalHeight = height;
		this.superSampling = superSampling;
		this.helmet = helmet;
		this.shadow = shadow;
		this.lighting = lighting;
		this.useWindow = useWindow;
	}

	public BufferedImage render(BufferedImage skin) throws Exception {
		init();
		BufferedImage head = skin.getSubimage(0, 0, 32, 16);
		BufferedImage helm = skin.getSubimage(32, 0, 32, 16);
		if (useWindow) {
			GL11.glClearColor(1, 1, 1, 1);
		}
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
		if (helmet) {
			upload(helm, 2);
		}

		if (shadow) {
			GL11.glEnable(GL11.GL_BLEND);
			GL11.glDisable(GL11.GL_TEXTURE_2D);
			GL11.glDisable(GL11.GL_LIGHTING);
			GL11.glPushMatrix();
				GL11.glTranslatef(0f, -0.95f, -0.45f);
				float scale = 1.02f;
				int count = 10;
				for (int i = 0; i < count; i++) {
					scale += 0.01f;
					GL11.glTranslatef(0f, -0.001f, 0f);
					GL11.glColor4f(0, 0, 0, (1-(i/(float)count))/2f);
					draw(scale, 0.01f, scale);
				}
			GL11.glPopMatrix();
		}

		GL11.glEnable(GL11.GL_TEXTURE_2D);
		if (lighting) {
			GL11.glEnable(GL11.GL_LIGHTING);
			GL11.glEnable(GL11.GL_LIGHT0);
		}

		GL11.glBlendFunc(GL11.GL_SRC_ALPHA, GL11.GL_ONE_MINUS_SRC_ALPHA);
		GL11.glColor3f(1, 1, 1);
		GL11.glBindTexture(GL11.GL_TEXTURE_2D, 1);
		draw(1.0f, 1.0f, 1.0f);
		if (helmet) {
			GL11.glBindTexture(GL11.GL_TEXTURE_2D, 2);
			draw(1.05f, 1.05f, 1.05f);
		}

		if (useWindow) {
			return null;
		} else {
			BufferedImage img = readPixels();
			BufferedImage out = new BufferedImage(finalWidth, finalHeight, BufferedImage.TYPE_INT_ARGB);
			Graphics2D gout = out.createGraphics();
			gout.drawImage(img.getScaledInstance(finalWidth, finalHeight, Image.SCALE_SMOOTH), 0, 0, null);
			gout.dispose();
			cleanup();
			return out;
		}
	}

	private BufferedImage readPixels() {
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
	private void draw(float xScale, float yScale, float zScale) throws Exception {
		GL11.glPushMatrix();
		GL11.glRotatef(20,1.0f,0f,0.0f);
		GL11.glTranslatef(0,-1.5f,-4.5f);
		GL11.glRotatef(angle,0f,1.0f,0f);
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

	private void createBuffer() throws Exception {
		buffer = new Pbuffer(width, height, new PixelFormat(), null, null);
		buffer.makeCurrent();
		if (buffer.isBufferLost()) {
			cleanup();
			System.err.println("Failed to set up PBuffer.");
			System.exit(2);
		}
	}
	
	private void createWindow() {
		try {
			Display.setFullscreen(false);
			Display.setTitle("Lapitar");
			Display.setDisplayMode(new DisplayMode(width, height));
			Display.create();
		} catch (LWJGLException e) {
			cleanup();
			System.err.println("Failed to set up Display.");
			System.exit(2);
		}
	}
	
	private void init() throws Exception {
		if (initialized) return;
		if (useWindow) {
			createWindow();
		} else {
			createBuffer();
		}
		initGL();
		initialized = true;
	}

	private void initGL() {
		GL11.glClearColor(0.0f, 0.0f, 0.0f, 0.0f);
		GL11.glClearDepth(1.0);
		GL11.glShadeModel(GL11.GL_SMOOTH);
		GL11.glEnable(GL11.GL_DEPTH_TEST);
		GL11.glDepthFunc(GL11.GL_LEQUAL);

		GL11.glMatrixMode(GL11.GL_PROJECTION);
		GL11.glLoadIdentity();

		GLU.gluPerspective(
				45.0f,
				(float) width / (float) height,
				0.1f,
				100.0f);
		GL11.glMatrixMode(GL11.GL_MODELVIEW);

		GL11.glHint(GL11.GL_PERSPECTIVE_CORRECTION_HINT, GL11.GL_NICEST);
	}

	private void cleanup() {
		if (useWindow) {
			Display.destroy();
		} else {
			buffer.destroy();
		}
	}

	@Override
	public String toString() {
		return MoreObjects.toStringHelper(this)
				.add("angle", angle)
				.add("width", finalWidth)
				.add("height", finalHeight)
				.add("superSampling", superSampling)
				.add("helmet", helmet)
				.add("shadow", shadow)
				.add("lighting", lighting)
				.toString();
	}

	public void incrementAngle() {
		angle += 0.5f;
	}
}
