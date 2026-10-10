import argparse
import hashlib
import re
import shutil
import zipfile
from pathlib import Path

root = Path(__file__).resolve().parents[1]
parser = argparse.ArgumentParser()
parser.add_argument('--abi', action='append', choices=['arm64-v8a', 'armeabi-v7a', 'x86_64'])
options = parser.parse_args()
match = re.search(
    r'^version:\s*([\w.+-]+)\s*$',
    (root / 'pubspec.yaml').read_text(encoding='utf-8'),
    re.MULTILINE,
)
if not match:
    raise SystemExit('pubspec.yaml 缺少合法版本号。')
version = match.group(1)
output = root / 'dist' / 'android6'
output.mkdir(parents=True, exist_ok=True)
artifacts = []
for abi in options.abi or ['arm64-v8a', 'armeabi-v7a', 'x86_64']:
    source = root / 'build' / 'app' / 'outputs' / 'flutter-apk' / f'app-{abi}-release.apk'
    if not source.is_file():
        raise SystemExit('缺少 Android 6 兼容版 APK：' + str(source))
    with zipfile.ZipFile(source) as archive:
        names = set(archive.namelist())
        required = [
            f'lib/{abi}/libduanju_core.so',
            f'lib/{abi}/libflutter.so',
            f'lib/{abi}/libapp.so',
            f'lib/{abi}/libmpv.so',
        ]
        missing = set(required) - names
        if missing:
            raise SystemExit('兼容版 APK 缺少原生库：' + ', '.join(sorted(missing)))
    target = output / f'hongguojian-android6-{version}-{abi}.apk'
    shutil.copy2(source, target)
    artifacts.append(target)

checksums = []
for artifact in sorted(artifacts):
    digest = hashlib.sha256()
    with artifact.open('rb') as stream:
        for chunk in iter(lambda: stream.read(1024 * 1024), b''):
            digest.update(chunk)
    checksums.append(f'{digest.hexdigest()}  {artifact.name}')
    print(artifact)
(output / 'SHA256SUMS.txt').write_text('\n'.join(checksums) + '\n', encoding='ascii')
