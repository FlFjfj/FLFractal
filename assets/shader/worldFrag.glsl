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


vec2 CM(vec2 a, vec2 b) {
    return vec2(a.x*b.x - a.y*b.y, a.x*b.y + a.y*b.x);
}



void main(){

    float r = 0.7885;
    vec2 c = vec2(r*cos(u_Delta*0.3), r*sin(u_Delta*0.3));
    float R = (1 + sqrt(1 + 4 * length(c)*length(c)))/2;
    vec2 z = v_Position.xy*0.3;

    vec2 c2 = vec2(r*cos(u_Delta*0.3 + 3.14), r*sin(u_Delta*0.3 + 3.14));
    float R2 = (1 + sqrt(1 + 4 * length(c2)*length(c2)))/2;
    vec2 z2 = v_Position.yx*0.3;


    float i = 0;
    float maxiter = 100;
    for(i = 0; i < maxiter; i += 1) {
        z = CM(z,z) + c;
        if( length(z) > R) {
            break;
        }
    }

     float i2 = 0;
     for(i2 = 0; i2 < maxiter; i2 += 1) {
             z2 = CM(z2,z2) + c2;
             if( length(z2) > R2) {
                 break;
             }
         }
    color = vec4(0.1 + i/maxiter,0.1 + i2/maxiter,(i*0.3+i2*0.8)/maxiter,1);
    vec3 lightPos = vec3(vec2(cos(u_Delta*0.3), sin(u_Delta*0.3)) * 2.5, 1);
	color = color * //vec4(1.0, 0.3, 0.7, 1.0) *
			(abs(dot(normalize(lightPos - v_Position),
					 normalize(v_VertNormal)))* 0.9 + 0.6);
}