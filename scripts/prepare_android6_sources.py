from pathlib import Path


root = Path(__file__).resolve().parents[1]
replacements = {
    root / 'lib' / 'home_screen.dart':
        '                            scrollCacheExtent: ScrollCacheExtent.viewport(1),\n',
    root / 'lib' / 'remote_widgets.dart':
        '      scrollCacheExtent: const ScrollCacheExtent.pixels(1200),\n',
}
for path, old in replacements.items():
    text = path.read_text(encoding='utf-8')
    if old not in text:
        raise SystemExit(f'兼容版源码缺少预期的 Flutter 新版参数：{path}')
    path.write_text(text.replace(old, ''), encoding='utf-8')
