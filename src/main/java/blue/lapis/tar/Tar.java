package blue.lapis.tar;

import static java.util.Arrays.asList;

import java.awt.image.BufferedImage;
import java.io.IOException;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;

import javax.imageio.ImageIO;

import joptsimple.OptionException;
import joptsimple.OptionParser;
import joptsimple.OptionSet;
import joptsimple.OptionSpec;

import org.lwjgl.input.Mouse;
import org.lwjgl.opengl.Display;

import com.google.common.base.Stopwatch;

public final class Tar {
	private static final float DEFAULT_ANGLE = 45;
	private static final float DEFAULT_TILT = 20;
	private static final int DEFAULT_SIZE = 256;

	public static void main(String[] args) {
		OptionParser parser = new OptionParser();
		parser.acceptsAll(asList("help", "?", "h"), "Show this help page.").forHelp();

		OptionSpec<SkinSource.Loader> loader = parser.acceptsAll(asList("d", "dl", "download"), "The method to get the skins from. Can be either MOJANG, URL or FILE.")
				.withRequiredArg().ofType(SkinSource.Loader.class).defaultsTo(SkinSource.Loader.MOJANG);
		OptionSpec<Output> output = parser.acceptsAll(asList("o", "out", "output"), "The destination to write the result to. Can be either FILE or STDOUT.")
				.withRequiredArg().ofType(Output.class).defaultsTo(Output.FILE);

		OptionSpec<Float> angle = parser.acceptsAll(asList("angle", "a"), "The angle to render the head at, in degrees.")
				.withRequiredArg().ofType(Float.class).defaultsTo(DEFAULT_ANGLE);
		OptionSpec<Float> tilt = parser.acceptsAll(asList("tilt", "t"), "The tilt to render the head at, in degrees.")
				.withRequiredArg().ofType(Float.class).defaultsTo(DEFAULT_TILT);

		OptionSpec<Integer> width = parser.acceptsAll(asList("width", "w"), "The width of the canvas to render on, in pixels.")
				.withRequiredArg().ofType(Integer.class).defaultsTo(DEFAULT_SIZE);
		OptionSpec<Integer> height = parser.acceptsAll(asList("height", "h"), "The height of the canvas to render on, in pixels.")
				.withRequiredArg().ofType(Integer.class).defaultsTo(DEFAULT_SIZE);

		OptionSpec<Integer> superSampling = parser.acceptsAll(asList("supersampling", "s"), "The amount of super sampling to perform, as a multiplier to width and height.")
				.withRequiredArg().ofType(Integer.class).defaultsTo(4);
		OptionSpec<Float> zoom = parser.acceptsAll(asList("zoom", "z"), "The amount to zoom in or out.")
				.withRequiredArg().ofType(Float.class).defaultsTo(-4.5f);
		
		parser.acceptsAll(asList("continuous"), "Continuously render a head in a window. Intended for debugging. Doesn't work properly with supersampling.");

		parser.accepts("no-helm", "Don't render the helm of the skin.");
		parser.accepts("no-shadow", "Don't render the shadow.");
		parser.accepts("no-lighting", "Don't enable lighting.");
		parser.accepts("portrait", "Render the head, torso, and arms instead of just the head.");
		parser.accepts("body", "Render the head, torso, arms, and legs instead of just the head.");
		parser.accepts("isometric", "Render in ugly isometric mode instead of perspective mode.");

		OptionSet options;
		try {
			options = parser.parse(args);
		} catch (OptionException e) {
			System.err.println(e.getMessage());
			System.exit(1);
			return;
		}

		if (options.has("help")) {
			try {
				parser.printHelpOn(System.err);
				System.err.println("Example: java -jar Tar.jar --angle 65 --tilt 30 Aesen_ Minecrell");
				System.err.println("\tGenerates two output files, Minecrell.png, and Aesen_.png, with the given arguments");
			} catch (IOException e) {
				e.printStackTrace();
				System.exit(2);
			}
		} else {
			boolean stdout = options.has("stdout");
			boolean continuous = options.has("continuous");
			Stopwatch watch = Stopwatch.createUnstarted();

			System.err.println("Initializing renderer...");
			watch.start();
			TarRenderer renderer = new TarRenderer(
					options.valueOf(angle),
					options.valueOf(tilt),
					options.valueOf(width),
					options.valueOf(height),
					options.valueOf(superSampling),
					!options.has("no-helm"),
					!options.has("no-shadow"),
					!options.has("no-lighting"),
					options.has("portrait") || options.has("body"),
					continuous,
					options.has("isometric"),
					options.has("body"),
					options.valueOf(zoom)
			);
			watch.stop();
			System.err.println("Initialized renderer in " + Time.format(watch));
			System.err.println("Rendering skins using " + renderer);
			watch.reset();

			System.err.println();

			System.err.println("Downloading skins...");
			SkinSource.Loader downloader = options.valueOf(loader);

			List<?> sources = options.nonOptionArguments();
			if (sources.size() > 1 && stdout) {
				System.err.println("Cannot bulk generate avatars when --stdout is specified!");
				System.exit(1);
			}

			Map<SkinSource, BufferedImage> skins = new LinkedHashMap<>(sources.size());

			long time = 0;
			for (Object arg : sources) {
				SkinSource source;
				try {
					source = downloader.create(arg.toString());
				} catch (IOException e) {
					System.err.println("Failed to initialize download for " + arg.toString());
					e.printStackTrace();
					continue;
				}

				watch.start();
				try {
					skins.put(source, ImageIO.read(source.getURL()));
				} catch (IOException e) {
					watch.stop();
					System.err.println("Failed to download skin for " + source + " (" + Time.format(watch) + ")");
					e.printStackTrace();
					watch.reset();
					continue;
				}
				watch.stop();
				System.err.println("Downloaded skin: " + source + " (" + Time.format(watch) + ")");
				time += watch.elapsed(TimeUnit.NANOSECONDS);
				watch.reset();
			}
			System.err.println("Finished downloading skins, took " + Time.format(time));

			System.err.println();
            System.err.println("Initializing LWJGL");
            try {
                Class.forName("org.lwjgl.opengl.Display");
            } catch (ClassNotFoundException e1) {
                e1.printStackTrace();
            }

			System.err.println("Rendering skins...");
			Map<SkinSource, BufferedImage> results = new LinkedHashMap<>(skins.size());

			boolean go = true;
			int fpsCounter = 0;
			int fps = 0;
			long lastFrameUpdate = 0;
			while (go) {
				time = 0;
				for (Map.Entry<SkinSource, BufferedImage> job : skins.entrySet()) {
					SkinSource source = job.getKey();
					BufferedImage skin = job.getValue();
	
					watch.start();
					try {
						results.put(source, renderer.render(skin));
					} catch (Exception e) {
						watch.stop();
						System.err.println("Failed to render skin for " + source + " (" + Time.format(watch) + ")");
						e.printStackTrace();
						watch.reset();
						continue;
					}
					watch.stop();
					if (!continuous) {
						System.err.println("Rendered skin: " + source + " (" + Time.format(watch) + ")");
					}
					time += watch.elapsed(TimeUnit.NANOSECONDS);
					watch.reset();
				}
				if (!continuous) {
					System.err.println("Finished rendering skins, took " + Time.format(time));
				}
				if (continuous) {
					while (!Display.isCreated()) { try { Thread.sleep(10L); } catch (InterruptedException e) {} }
					if (Display.isCloseRequested()) {
						go = false;
					}
					fpsCounter++;
					if (System.currentTimeMillis() - lastFrameUpdate >= 1000) {
						fps = fpsCounter;
						fpsCounter = 0;
						lastFrameUpdate = System.currentTimeMillis();
					}
					renderer.modifyZoom(Mouse.getDWheel()/240f);
					if (Mouse.isButtonDown(1)) {
						renderer.modifyAngle(Mouse.getDX(), -Mouse.getDY());
					} else {
						Mouse.getDX();
						Mouse.getDY();
					}
					Display.setTitle("Lapitar - "+fps+"/30 fps");
					Display.update();
					Display.sync(30);
				} else {
					go = false;
				}
			}

			System.err.println();

			if (!continuous) {
				System.err.println("Saving results...");
				Output out = options.valueOf(output);
				time = 0;
				for (Map.Entry<SkinSource, BufferedImage> result : results.entrySet()) {
					SkinSource source = result.getKey();
					BufferedImage skin = result.getValue();
	
					watch.start();
					try {
						out.write(source, skin);
					} catch (IOException e) {
						watch.stop();
						System.err.println("Failed to save result for " + source + " (" + Time.format(watch) + ")");
						e.printStackTrace();
						watch.reset();
						continue;
					}
					watch.stop();
					System.err.println("Saved result: " + source + " (" + Time.format(watch) + ")");
					time += watch.elapsed(TimeUnit.NANOSECONDS);
					watch.reset();
				}
				System.err.println("Finished saving results, took " + Time.format(time));
			}
			
			System.err.println();

			System.err.println("Done!");
		}
	}
}
