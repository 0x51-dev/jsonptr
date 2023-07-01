# ABNF

```abnf
unescaped             = %x00-2E / %x30-7D / %x7F-10FFFF
--- becomes ---
unescaped             = %x00-2E / %x30-7D / %x7F-FFFC / %xFFFE-10FFFF
```