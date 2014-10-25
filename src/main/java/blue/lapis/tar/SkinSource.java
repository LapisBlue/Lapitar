package blue.lapis.tar;

import com.google.common.io.Files;

import java.io.IOException;
import java.net.URL;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Objects;

public class SkinSource {
	private final String player;
	private final URL url;

	public SkinSource(String player, URL url) {
		this.player = Objects.requireNonNull(player, "player");
		this.url = Objects.requireNonNull(url, "url");
	}

	public String getPlayer() {
		return player;
	}

	public URL getURL() {
		return url;
	}

	@Override
	public boolean equals(Object obj) {
		if (this == obj) return true;
		if (obj == null || !(obj instanceof SkinSource)) return false;
		SkinSource other = (SkinSource) obj;
		return Objects.equals(player, other.player)
				&& Objects.equals(url, other.url);
	}

	@Override
	public int hashCode() {
		return Objects.hash(player, url);
	}

	@Override
	public String toString() {
		return player;
	}

	public enum Loader {
		MOJANG {
			@Override
			public SkinSource create(String player) throws IOException {
				return new SkinSource(player, new URL(String.format("http://skins.minecraft.net/MinecraftSkins/%s.png", player)));
			}
		}, URL {
			@Override
			public SkinSource create(String input) throws IOException {
				String player = input;
				int pos = input.lastIndexOf('/');
				if (pos != -1)
					player = input.substring(pos + 1);
				return new SkinSource(normalizeFile(player), new URL(input));
			}
		}, FILE {
			@Override
			public SkinSource create(String input) throws IOException {
				Path path = Paths.get(input);
				return new SkinSource(normalizeFile(path.getFileName().toString()), path.toUri().toURL());
			}
		};

		public abstract SkinSource create(String input) throws IOException;

		private static String normalizeFile(String file) {
			return Files.getNameWithoutExtension(file);
		}
	}
}
