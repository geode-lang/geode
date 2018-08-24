is vector


class Vec2 {
	float x;
	float y;
}

func NewVec2(float x, float y) Vec2 {
	Vec2 v;
	v.x <- x;
	v.y <- y;
	return v;
}

func vec2_dot(Vec2 a, Vec2 b) float {
	return a.x * b.x + a.y * b.y;
}