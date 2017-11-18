#version 330 core

uniform mat4 u_ProjTrans;
uniform mat4 u_ObjTrans;

uniform vec3 u_Color;

layout (location = 0) in vec4 a_Position;
layout (location = 1) in vec2 a_TexCoord;
layout (location = 2) in vec3 a_VertNormal;

out vec2 v_TexCoord;
out vec3 v_VertNormal;
out vec3 v_Position;

void main() {
   v_TexCoord = a_TexCoord;
   v_Position = (u_ObjTrans * a_Position).xyz;
   v_VertNormal = vec3(0.0, 0.0, 1.0);
   gl_Position = u_ProjTrans * u_ObjTrans * a_Position;
}