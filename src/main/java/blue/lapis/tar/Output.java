package blue.lapis.tar;

import javax.imageio.ImageIO;
import java.awt.image.BufferedImage;
import java.io.IOException;
import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Paths;

public enum Output {
	STDOUT {
		@Override
		public void write(SkinSource source, BufferedImage result) throws IOException {
			ImageIO.write(result, "png", System.out);
		}
	}, FILE {
		@Override
		public void write(SkinSource source, BufferedImage result) throws IOException {
			try (OutputStream out = Files.newOutputStream(Paths.get(source.getPlayer() + ".png"))) {
				ImageIO.write(result, "png", out);
			}
		}
	};

	public abstract void write(SkinSource source, BufferedImage result) throws IOException;
}
