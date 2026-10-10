import argparse
import os
import platform
import shutil
import subprocess
import sys
from pathlib import Path

root = Path(__file__).resolve().parents[1]
parser = argparse.ArgumentParser()
parser.add_argument('--abi', action='append', choices=['arm64-v8a', 'armeabi-v7a', 'x86_64'])
options = parser.parse_args()
environment = os.environ.copy()
if platform.system() == 'Darwin':
    environment['LANG'] = 'en_US.UTF-8'
    environment['LC_ALL'] = 'en_US.UTF-8'
environment.setdefault('GOPROXY', 'https://proxy.golang.org,direct')
environment.setdefault('GOSUMDB', 'sum.golang.org')
flutter = shutil.which('flutter')
if not flutter:
    raise SystemExit('请先将 Flutter SDK 的 bin 目录加入 PATH。')
abi_args = [item for abi in options.abi or [] for item in ['--abi', abi]]
subprocess.run(
    [
        sys.executable,
        str(root / 'scripts' / 'build_native.py'),
        '--platform',
        'android',
        '--android-api',
        '23',
        *abi_args,
    ],
    cwd=root,
    env=environment,
    check=True,
)
subprocess.run([flutter, 'pub', 'get'], cwd=root, env=environment, check=True)
build_args = [
    flutter,
    'build',
    'apk',
    '--release',
    '--split-per-abi',
    '--no-pub',
    '--dart-define=ALL_SOURCES=false',
    '--dart-define=ANDROID6_COMPAT=true',
]
if options.abi:
    targets = {
        'arm64-v8a': 'android-arm64',
        'armeabi-v7a': 'android-arm',
        'x86_64': 'android-x64',
    }
    build_args += ['--target-platform', ','.join(targets[abi] for abi in options.abi)]
subprocess.run(build_args, cwd=root, env=environment, check=True)
subprocess.run(
    [sys.executable, str(root / 'scripts' / 'package_android6.py'), *abi_args],
    cwd=root,
    env=environment,
    check=True,
)
