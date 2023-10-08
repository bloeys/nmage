//shader:vertex
#version 410

uniform mat4 ProjMtx;

in vec2 Position;
in vec2 UV;
in vec4 Color;

out vec2 Frag_UV;
out vec4 Frag_Color;

// Imgui doesn't handle srgb correctly, and looks too bright and wrong in srgb buffers (see: https://github.com/ocornut/imgui/issues/578).
// While not a complete fix (that would require changes in imgui itself), moving incoming srgba colors to linear in the vertex shader helps make things look better.
vec4 srgba_to_linear(vec4 srgbaColor){

    #define gamma_correction 2.2

    return vec4(
        pow(srgbaColor.r, gamma_correction),
        pow(srgbaColor.g, gamma_correction),
        pow(srgbaColor.b, gamma_correction),
        srgbaColor.a
    );
}

void main()
{
    Frag_UV = UV;
    Frag_Color = srgba_to_linear(Color);
    gl_Position = ProjMtx * vec4(Position.xy, 0, 1);
}

//shader:fragment
#version 410

uniform sampler2D Texture;

in vec2 Frag_UV;
in vec4 Frag_Color;

out vec4 Out_Color;

void main()
{
    Out_Color = vec4(Frag_Color.rgb, Frag_Color.a * texture(Texture, Frag_UV.st).r);
}