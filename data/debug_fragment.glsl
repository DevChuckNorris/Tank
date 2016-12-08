#version 330

// Output data
out vec4 color;

uniform sampler2D tex;

in vec2 UV;

void main(){
	color = texture(tex, UV);
}
