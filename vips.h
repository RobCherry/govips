#include <stdbool.h>
#include <stdlib.h>
#include <vips/vips.h>

int govips_gifload_buffer(void *input, size_t length, VipsImage **output, gint page, VipsAccess access, gboolean disc) {
  return vips_gifload_buffer(input, length, output,
    "page", page,
    "access", access,
    "disc", disc,
    NULL);
}

int govips_jpegload_buffer(void *input, size_t length, VipsImage **output, gint shrink, gboolean fail, gboolean autorotate, VipsAccess access, gboolean disc) {
  return vips_jpegload_buffer(input, length, output,
    "shrink", shrink,
    "fail", fail,
    "autorotate", autorotate,
    "access", access,
    "disc", disc,
    NULL);
}

int govips_jpegsave_buffer(VipsImage *input, void **output, size_t *length, gint Q, const char *profile, gboolean optimize_coding, gboolean interlace, gboolean strip, gboolean nosubsample, gboolean trellis_quant, gboolean overshoot_deringing, gboolean optimize_scans, gint quant_table) {
  return vips_jpegsave_buffer(input, output, length,
    "Q", Q,
    "profile", profile,
    "optimize_coding", optimize_coding,
    "interlace", interlace,
    "strip", strip,
    "no-subsample", nosubsample,
    "trellis_quant", trellis_quant,
    "overshoot_deringing", overshoot_deringing,
    "optimize_scans", optimize_scans,
    "quant_table", quant_table,
    NULL);
}

int govips_jpegsave(VipsImage *input, const char *filename, gint Q, const char *profile, gboolean optimize_coding, gboolean interlace, gboolean strip, gboolean nosubsample, gboolean trellis_quant, gboolean overshoot_deringing, gboolean optimize_scans, gint quant_table) {
  return vips_jpegsave(input, filename,
    "Q", Q,
    "profile", profile,
    "optimize_coding", optimize_coding,
    "interlace", interlace,
    "strip", strip,
    "no-subsample", nosubsample,
    "trellis_quant", trellis_quant,
    "overshoot_deringing", overshoot_deringing,
    "optimize_scans", optimize_scans,
    "quant_table", quant_table,
    NULL);
}

int govips_magickload_buffer(void *input, size_t length, VipsImage **output, gboolean all_frames, const char *density, gint page, VipsAccess access, gboolean disc) {
  return vips_magickload_buffer(input, length, output,
    "all_frames", all_frames,
    "density", density,
    "page", page,
    "access", access,
    "disc", disc,
    NULL);
}

int govips_pngload_buffer(void *input, size_t length, VipsImage **output, VipsAccess access, gboolean disc) {
  return vips_pngload_buffer(input, length, output,
    "access", access,
    "disc", disc,
    NULL);
}

int govips_pngsave_buffer(VipsImage *input, void **output, size_t *length, gint compression, gboolean interlace, const char *profile, VipsForeignPngFilter filter) {
  return vips_pngsave_buffer(input, output, length,
    "compression", compression,
    "interlace", interlace,
    "profile", profile,
    "filter", filter,
    NULL);
}

int govips_pngsave(VipsImage *input, const char *filename, gint compression, gboolean interlace, const char *profile, VipsForeignPngFilter filter) {
  return vips_pngsave(input, filename,
    "compression", compression,
    "interlace", interlace,
    "profile", profile,
    "filter", filter,
    NULL);
}

int govips_webpload_buffer(void *input, size_t length, VipsImage **output, gint shrink, VipsAccess access, gboolean disc) {
  return vips_webpload_buffer(input, length, output,
    "shrink", shrink,
    "access", access,
    "disc", disc,
    NULL);
}

int govips_webpsave_buffer(VipsImage *input, void **output, size_t *length, gint Q, gboolean lossless, VipsForeignWebpPreset preset, gboolean smart_subsample, gboolean near_lossless, gint alpha_q) {
  return vips_webpsave_buffer(input, output, length,
    "Q", Q,
    "lossless", lossless,
    "preset", preset,
    "smart_subsample", smart_subsample,
    "near_lossless", near_lossless,
    "alpha_q", alpha_q,
    "strip", 1,
    NULL);
}

int govips_webpsave(VipsImage *input, const char *filename, gint Q, gboolean lossless, VipsForeignWebpPreset preset, gboolean smart_subsample, gboolean near_lossless, gint alpha_q) {
  return vips_webpsave(input, filename,
    "Q", Q,
    "lossless", lossless,
    "preset", preset,
    "smart_subsample", smart_subsample,
    "near_lossless", near_lossless,
    "alpha_q", alpha_q,
    "strip", 1,
    NULL);
}

int govips_embed(VipsImage *in, VipsImage **out, int x, int y, int width, int height, VipsExtend extend, VipsArrayDouble *background) {
  if (extend == VIPS_EXTEND_BACKGROUND && background != NULL) {
    return vips_embed(in, out, x, y, width, height, "extend", extend, "background", background, NULL);
  }
  return vips_embed(in, out, x, y, width, height, "extend", extend, NULL);
}

int govips_extract_area(VipsImage *in, VipsImage **out, int left, int top, int width, int height) {
  return vips_extract_area(in, out, left, top, width, height, NULL);
}

int govips_shrink(VipsImage *in, VipsImage **out, double xshrink, double yshrink) {
  return vips_shrink(in, out, xshrink, yshrink, NULL);
}

int govips_shrinkh(VipsImage *in, VipsImage **out, double xshrink) {
  return vips_shrinkh(in, out, xshrink, NULL);
}

int govips_shrinkv(VipsImage *in, VipsImage **out, double yshrink) {
  return vips_shrinkv(in, out, yshrink, NULL);
}

int govips_reduce(VipsImage *in, VipsImage **out, double xshrink, double yshrink, VipsKernel kernel) {
  return vips_reduce(in, out, xshrink, yshrink, "kernel", kernel, NULL);
}

int govips_reduceh(VipsImage *in, VipsImage **out, double xshrink, VipsKernel kernel) {
  return vips_reduceh(in, out, xshrink, "kernel", kernel, NULL);
}

int govips_reducev(VipsImage *in, VipsImage **out, double yshrink, VipsKernel kernel) {
  return vips_reducev(in, out, yshrink, "kernel", kernel, NULL);
}

int govips_resize(VipsImage *in, VipsImage **out, double scale, double vscale, VipsKernel kernel) {
  return vips_resize(in, out, scale, "vscale", vscale, "kernel", kernel, NULL);
}

int govips_similarity(VipsImage *in, VipsImage **out, gdouble scale, gdouble angle, VipsInterpolate *interpolate, gdouble idx, gdouble idy, gdouble odx, gdouble ody) {
  return vips_similarity(in, out, "scale", scale, "angle", angle, "interpolate", interpolate, "idx", idx, "idy", idy, "odx", odx, "ody", ody, NULL);
}

int govips_affine(VipsImage *in, VipsImage **out, double a, double b, double c, double d, VipsInterpolate *interpolate, VipsArrayInt *oarea, gdouble idx, gdouble idy, gdouble odx, gdouble ody) {
  if (oarea == NULL) {
    return vips_affine(in, out, a, b, c, d, "interpolate", interpolate, "idx", idx, "idy", idy, "odx", odx, "ody", ody, NULL);
  }
  return vips_affine(in, out, a, b, c, d, "interpolate", interpolate, "oarea", oarea, "idx", idx, "idy", idy, "odx", odx, "ody", ody, NULL);
}

int govips_gaussblur(VipsImage *in, VipsImage **out, double sigma, VipsPrecision precision, double min_ampl) {
  return vips_gaussblur(in, out, sigma, "precision", precision, "min_ampl", min_ampl, NULL);
}

int govips_sharpen(VipsImage *in, VipsImage **out, double sigma, double x1, double y2, double y3, double m1, double m2) {
  return vips_sharpen(in, out, "sigma", sigma, "x1", x1, "y2", y2, "y3", y3, "m1", m1, "m2", m2, NULL);
}

int govips_flatten(VipsImage *in, VipsImage **out, VipsArrayDouble *background, double max_alpha) {
  if (background == NULL) {
    return vips_flatten(in, out, "max_alpha", max_alpha, NULL);
  }
  return vips_flatten(in, out, "background", background, "max_alpha", max_alpha, NULL);
}

int govips_colourspace(VipsImage *in, VipsImage **out, VipsInterpretation space, VipsInterpretation source_space) {
  return vips_colourspace(in, out, space, "source_space", source_space, NULL);
}

int govips_icc_transform(VipsImage *in, VipsImage **out, const char *output_profile, const char *input_profile, VipsIntent intent, int depth, gboolean embedded) {
  return vips_icc_transform(in, out, output_profile, "input_profile", input_profile, "intent", intent, "depth", depth, "embedded", embedded, NULL);
}

VipsRect govips_rect_new(int left, int top, int width, int height) {
  VipsRect r = { .left = left, .top = top, .width = width, .height = height };
  return r;
}

VipsPel* govips_region_addr(VipsRegion *r, int x, int y) {
  return VIPS_REGION_ADDR(r, x, y);
}

int govips_vips_region_n_elements(VipsRegion *r) {
  return VIPS_REGION_N_ELEMENTS(r);
}

VipsPel* govips_pel_band(VipsPel *p, int n) {
  return p + n;
}
