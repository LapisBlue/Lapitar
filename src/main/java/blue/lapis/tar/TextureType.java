package blue.lapis.tar;



public enum TextureType {
	NONE,
	HEAD(
		32, 16,
		// Front (Red)
		8, 8, 8, 8,
		// Back (Blue)
		24, 8, 8, 8,
		// Top (Purple)
		8, 0, 8, 8,
		// Bottom (Gray)
		16, 0, 8, 8,
		// Left (Yellow)
		16, 8, 8, 8,
		// Right (Green)
		0, 8, 8, 8
	),
	TORSO(
		24, 16,
		// Front (Red)
		4, 4, 8, 12,
		// Back (Blue)
		16, 4, 8, 12,
		// Top (Purple)
		4, 0, 8, 4,
		// Bottom (Gray)
		12, 0, 8, 4,
		// Left (Yellow)
		12, 4, 4, 12,
		// Right (Green)
		0, 4, 4, 12
	),
	ARM(
		16, 16,
		// Front (Red)
		4, 4, 4, 12,
		// Back (Blue)
		12, 4, 4, 12,
		// Top (Purple)
		4, 0, 4, 4,
		// Bottom (Gray)
		8, 0, 4, 4,
		// Left (Yellow)
		8, 4, 4, 12,
		// Right (Green)
		0, 4, 4, 12
	);
	
	public final float[] u = new float[24];
	public final float[] v = new float[24];
	
	private TextureType() {}
	
	// constructor uses varargs for compactness
	// arguments are effectively:
	// texture_width, texture_height
	// <side>_x, <side>_y, <side>_width, <side>_height
	// where <side> is the face in question, in order:
	// Front, Back, Top, Bottom, Left, Right
	private TextureType(int tex_w, int tex_h, int... assorted) {
		for (int i = 0; i < assorted.length/4; i++) {
			int idx = i*4;
			
			int x = assorted[idx];
			int y = assorted[idx+1];
			int edgeX = x + assorted[idx+2];
			int edgeY = y + assorted[idx+3];
			
			u[idx  ] = div(tex_w,     x);
			v[idx  ] = div(tex_h, edgeY);
			
			u[idx+1] = div(tex_w, edgeX);
			v[idx+1] = div(tex_h, edgeY);
			
			u[idx+2] = div(tex_w, edgeX);
			v[idx+2] = div(tex_h,     y);
			
			u[idx+3] = div(tex_w,     x);
			v[idx+3] = div(tex_h,     y);
		}
	}

	private float div(float max, float x) { // upcasting!
		return x / max;
	}

}
