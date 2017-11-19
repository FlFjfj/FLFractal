#version 330 core

uniform sampler2D tex;
uniform mat4 u_ProjTrans;
uniform mat4 u_ObjTrans;

uniform vec3 u_LightPos;
uniform float u_Delta;


in vec2 v_TexCoord;
in vec3 v_VertNormal;
in vec3 v_Position;
out vec4 color;

void main(){
    vec3 lightPos = vec3(vec2(cos(u_Delta), sin(u_Delta)) * 2.5, 1);
	color = vec4(1.0, 0.3, 0.7, 1.0) *
			(abs(dot(normalize(lightPos - v_Position),
					 normalize(v_VertNormal)))* 0.7 + 0.3);
}