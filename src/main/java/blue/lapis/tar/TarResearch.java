package blue.lapis.tar;

import java.awt.Graphics2D;
import java.awt.Image;
import java.awt.image.BufferedImage;
import java.io.File;
import java.io.IOException;
import java.net.URL;
import java.nio.ByteBuffer;
import java.nio.FloatBuffer;
import java.nio.IntBuffer;
import java.util.List;
import java.util.Map;

import javax.imageio.ImageIO;

import joptsimple.HelpFormatter;
import joptsimple.OptionDescriptor;
import joptsimple.OptionParser;
import joptsimple.OptionSet;
import joptsimple.OptionSpec;

import org.lwjgl.BufferUtils;
import org.lwjgl.opengl.GL11;
import org.lwjgl.opengl.Pbuffer;
import org.lwjgl.opengl.PixelFormat;
import org.lwjgl.util.glu.GLU;

public class TarResearch {
	public static int supersampling = 4;
	public static int width = 256*supersampling;
	public static int height = 256*supersampling;
	private static Pbuffer buffer;
	private static FloatBuffer lightPosition;
	private static FloatBuffer lightAmbient;
	
	public static void main(String[] args) {
		OptionParser parser = new OptionParser();
		parser.accepts("help", "Display this help.").forHelp();
		OptionSpec<Float> angle = parser.accepts("angle", "The angle to render the head at, in degrees.").withRequiredArg().ofType(Float.class).defaultsTo(45f);
		OptionSpec<Integer> widthSpec = parser.accepts("width", "The width of the canvas to render on, in pixels.").withRequiredArg().ofType(Integer.class).defaultsTo(256);
		OptionSpec<Integer> heightSpec = parser.accepts("height", "The height of the canvas to render on, in pixels.").withRequiredArg().ofType(Integer.class).defaultsTo(256);
		OptionSpec<Integer> supersamplingSpec = parser.accepts("supersampling", "The amount of supersampling to perform, as a multiplier to width and height").withRequiredArg().ofType(Integer.class).defaultsTo(4);
		parser.accepts("no-helm", "Don't render the helm portion of the skin.");
		parser.accepts("no-shadow", "Don't render the shadow.");
		parser.accepts("no-lighting", "Don't enable lighting.");
		OptionSet options = parser.parse(args);
		if (options.has("help")) {
			try {
				parser.printHelpOn(System.out);
				System.out.println("Example: java -jar Tar.jar --angle 65 --tilt 30 Aesen_ Minecrell");
				System.out.println("\tGenerates two output files, Minecrell.png, and Aesen_.png, with the given arguments");
			} catch (IOException e) {
				e.printStackTrace();
			}
			return;
		}
		quadRot = options.valueOf(angle);
		supersampling = options.valueOf(supersamplingSpec);
		width = options.valueOf(widthSpec)*supersampling;
		height = options.valueOf(heightSpec)*supersampling;
		boolean helmet = !options.has("no-helm");
		boolean shadow = !options.has("no-shadow");
		boolean lighting = !options.has("no-lighting");
		for (Object o : options.nonOptionArguments()) {
			String player = String.valueOf(o);
			System.out.println("Creating avatar for "+player+" ("+(width/supersampling)+"x"+(height/supersampling)+", "+supersampling+"x supersampling, "+(!shadow?"without":"with")+" shadow, "+(!helmet?"without":"with")+" helmet, "+(!lighting?"without":"with")+" lighting, angle "+quadRot+"\u00B0)");
			try {
				long startTime = System.currentTimeMillis();
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
				if (helmet) {
					upload(helm, 2);
				}
				
				if (shadow) {
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
				System.out.println("Render complete (took "+((System.currentTimeMillis()-startTime)/1000f)+"s), writing output file");
				startTime = System.currentTimeMillis();
				BufferedImage img = readPixels();
				BufferedImage out = new BufferedImage(width/supersampling, height/supersampling, BufferedImage.TYPE_INT_ARGB);
				Graphics2D gout = out.createGraphics();
				gout.drawImage(img.getScaledInstance(img.getWidth()/supersampling, img.getHeight()/supersampling, Image.SCALE_SMOOTH), 0, 0, null);
				gout.dispose();
				File file = new File(player+".png");
				ImageIO.write(out, "png", file);
				cleanup();
				System.out.println("Successfully created "+file+" (took "+((System.currentTimeMillis()-startTime)/1000f)+"s); "+(file.length()/1024f)+"KiB)");
			} catch (Exception e) {
				e.printStackTrace();
				System.err.println("Failed to render head for "+player);
				continue;
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
	private static float quadRot = 45;
	private static void draw(float xScale, float yScale, float zScale) throws Exception {
		GL11.glPushMatrix();
			GL11.glRotatef(20,1.0f,0f,0.0f);
			GL11.glTranslatef(0,-1.5f,-4.5f);
			GL11.glRotatef(quadRot,0f,1.0f,0f);
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
