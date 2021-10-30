#version 460

in vec3 vertPos;
in vec3 vertColor;

out vec3 outColor;

void main()
{
    outColor = vertColor;
    gl_Position = vec4(vertPos, 1.0); // vec4(vertPos.x, vertPos.y, vertPos.z, 1.0)
}