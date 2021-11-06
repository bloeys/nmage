#version 460

in vec3 vertPos;
in vec3 vertColor;

out vec3 outColor;

uniform mat4 modelMat;

void main()
{
    outColor = vertColor;
    gl_Position = modelMat * (vec4(vertPos, 1.0)); // vec4(vertPos.x, vertPos.y, vertPos.z, 1.0)
}