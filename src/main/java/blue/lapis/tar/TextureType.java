package blue.lapis.tar;


public enum TextureType {
	NONE(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0),
	HEAD(
		// Front
			0.25f, 1.00f,
			0.50f, 1.00f,
			0.50f, 0.50f,
			0.25f, 0.50f,
		// Back
			1.00f, 1.00f,
			1.00f, 0.50f,
			0.75f, 0.50f,
			0.75f, 1.00f,
		// Top
			0.50f, 0.00f,
			0.50f, 0.50f,
			0.25f, 0.50f,
			0.25f, 0.00f,
		// Bottom
			0.50f, 0.50f,
			0.75f, 0.50f,
			0.75f, 0.00f,
			0.50f, 0.00f,
		// Left
			0.75f, 1.00f,
			0.75f, 0.50f,
			0.50f, 0.50f,
			0.50f, 1.00f,
		// Right
			0.00f, 1.00f,
			0.25f, 1.00f,
			0.25f, 0.50f,
			0.00f, 0.50f
		),
	TORSO(
		// Front
			0.166666667f, 1.00f,
			0.50f, 1.00f,
			0.50f, 0.25f,
			0.166666667f, 0.25f,
		// Back
			0.666666667f, 1.00f,
			1.00f, 1.00f,
			1.00f, 0.25f,
			0.666666667f, 0.25f,
		// Top
			0.166666667f, 0.25f,
			0.50f, 0.25f,
			0.50f, 0.00f,
			0.166666667f, 0.00f,
		// Bottom
			0.50f, 0.25f,
			0.833333333f, 0.25f,
			0.833333333f, 0.00f,
			0.50f, 0.00f,
		// Left
			0.50f, 1.00f,
			0.666666667f, 1.00f,
			0.666666667f, 0.25f,
			0.50f, 0.25f,
		// Right
			0.00f, 1.00f,
			0.166666667f, 1.00f,
			0.166666667f, 0.25f,
			0.00f, 0.25f
		),
	ARM(
		// Front
			0.25f, 1.00f,
			0.50f, 1.00f,
			0.50f, 0.50f,
			0.25f, 0.50f,
		// Back
			1.00f, 1.00f,
			1.00f, 0.50f,
			0.75f, 0.50f,
			0.75f, 1.00f,
		// Top
			0.50f, 0.00f,
			0.50f, 0.50f,
			0.25f, 0.50f,
			0.25f, 0.00f,
		// Bottom
			0.50f, 0.50f,
			0.75f, 0.50f,
			0.75f, 0.00f,
			0.50f, 0.00f,
		// Left
			0.75f, 1.00f,
			0.75f, 0.50f,
			0.50f, 0.50f,
			0.50f, 1.00f,
		// Right
			0.00f, 1.00f,
			0.25f, 1.00f,
			0.25f, 0.50f,
			0.00f, 0.50f
		);
	
	// Front
	public final float fr_u_a;
	public final float fr_v_a;
	public final float fr_u_b;
	public final float fr_v_b;
	public final float fr_u_c;
	public final float fr_v_c;
	public final float fr_u_d;
	public final float fr_v_d;
	                         
	// Back                  
	public final float ba_u_a;
	public final float ba_v_a;
	public final float ba_u_b;
	public final float ba_v_b;
	public final float ba_u_c;
	public final float ba_v_c;
	public final float ba_u_d;
	public final float ba_v_d;
                             
	// Top                   
	public final float to_u_a;
	public final float to_v_a;
	public final float to_u_b;
	public final float to_v_b;
	public final float to_u_c;
	public final float to_v_c;
	public final float to_u_d;
	public final float to_v_d;
                             
	// Bottom                
	public final float bo_u_a;
	public final float bo_v_a;
	public final float bo_u_b;
	public final float bo_v_b;
	public final float bo_u_c;
	public final float bo_v_c;
	public final float bo_u_d;
	public final float bo_v_d;
	                         
	// Left                  
	public final float le_u_a;
	public final float le_v_a;
	public final float le_u_b;
	public final float le_v_b;
	
	// I'm so sorry.
	private TextureType(
			float fr_u_a, float fr_v_a, float fr_u_b, float fr_v_b,
			float fr_u_c, float fr_v_c, float fr_u_d, float fr_v_d,
			float ba_u_a, float ba_v_a, float ba_u_b, float ba_v_b,
			float ba_u_c, float ba_v_c, float ba_u_d, float ba_v_d,
			float to_u_a, float to_v_a, float to_u_b, float to_v_b,
			float to_u_c, float to_v_c, float to_u_d, float to_v_d,
			float bo_u_a, float bo_v_a, float bo_u_b, float bo_v_b,
			float bo_u_c, float bo_v_c, float bo_u_d, float bo_v_d,
			float le_u_a, float le_v_a, float le_u_b, float le_v_b,
			float le_u_c, float le_v_c, float le_u_d, float le_v_d,
			float ri_u_a, float ri_v_a, float ri_u_b, float ri_v_b,
			float ri_u_c, float ri_v_c, float ri_u_d, float ri_v_d) {
		this.fr_u_a = fr_u_a;
		this.fr_v_a = fr_v_a;
		this.fr_u_b = fr_u_b;
		this.fr_v_b = fr_v_b;
		this.fr_u_c = fr_u_c;
		this.fr_v_c = fr_v_c;
		this.fr_u_d = fr_u_d;
		this.fr_v_d = fr_v_d;
		this.ba_u_a = ba_u_a;
		this.ba_v_a = ba_v_a;
		this.ba_u_b = ba_u_b;
		this.ba_v_b = ba_v_b;
		this.ba_u_c = ba_u_c;
		this.ba_v_c = ba_v_c;
		this.ba_u_d = ba_u_d;
		this.ba_v_d = ba_v_d;
		this.to_u_a = to_u_a;
		this.to_v_a = to_v_a;
		this.to_u_b = to_u_b;
		this.to_v_b = to_v_b;
		this.to_u_c = to_u_c;
		this.to_v_c = to_v_c;
		this.to_u_d = to_u_d;
		this.to_v_d = to_v_d;
		this.bo_u_a = bo_u_a;
		this.bo_v_a = bo_v_a;
		this.bo_u_b = bo_u_b;
		this.bo_v_b = bo_v_b;
		this.bo_u_c = bo_u_c;
		this.bo_v_c = bo_v_c;
		this.bo_u_d = bo_u_d;
		this.bo_v_d = bo_v_d;
		this.le_u_a = le_u_a;
		this.le_v_a = le_v_a;
		this.le_u_b = le_u_b;
		this.le_v_b = le_v_b;
		this.le_u_c = le_u_c;
		this.le_v_c = le_v_c;
		this.le_u_d = le_u_d;
		this.le_v_d = le_v_d;
		this.ri_u_a = ri_u_a;
		this.ri_v_a = ri_v_a;
		this.ri_u_b = ri_u_b;
		this.ri_v_b = ri_v_b;
		this.ri_u_c = ri_u_c;
		this.ri_v_c = ri_v_c;
		this.ri_u_d = ri_u_d;
		this.ri_v_d = ri_v_d;
	}
	public final float le_u_c;
	public final float le_v_c;
	public final float le_u_d;
	public final float le_v_d;
                             
	// Right                 
	public final float ri_u_a;
	public final float ri_v_a;
	public final float ri_u_b;
	public final float ri_v_b;
	public final float ri_u_c;
	public final float ri_v_c;
	public final float ri_u_d;
	public final float ri_v_d;
	
	
}
