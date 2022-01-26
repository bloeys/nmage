#version 410

layout(location=0) in vec3 vertPosIn;
layout(location=1) in vec3 vertNormalIn;
layout(location=2) in vec3 vertColorIn;

out vec3 vertColor;
out vec3 vertNormal;
out vec3 fragPos;

//MVP = Model View Projection
uniform mat4 modelMat;
uniform mat4 viewMat;
uniform mat4 projMat;

void main()
{
    vertColor = vertColorIn;
    vertNormal = mat3(transpose(inverse(modelMat))) * vertNormalIn;
    fragPos = vec3(modelMat * vec4(vertPosIn, 1.0));

    gl_Position = projMat * viewMat * modelMat * vec4(vertPosIn, 1.0);
}