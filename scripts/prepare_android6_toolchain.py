from pathlib import Path


root = Path(__file__).resolve().parents[1]
app = root / 'android' / 'app' / 'build.gradle.kts'
app_text = app.read_text(encoding='utf-8').replace(
    '    id("com.android.application")\n    id("dev.flutter.flutter-gradle-plugin")',
    '    id("com.android.application")\n'
    '    id("org.jetbrains.kotlin.android")\n'
    '    id("dev.flutter.flutter-gradle-plugin")',
)
app.write_text(app_text, encoding='utf-8')

settings = root / 'android' / 'settings.gradle.kts'
settings_text = settings.read_text(encoding='utf-8')
settings_text = settings_text.replace(
    'id("com.android.application") version "9.1.0" apply false',
    'id("com.android.application") version "8.7.3" apply false',
).replace(
    'id("org.jetbrains.kotlin.android") version "2.4.0" apply false',
    'id("org.jetbrains.kotlin.android") version "2.1.0" apply false',
)
settings.write_text(settings_text, encoding='utf-8')

wrapper = root / 'android' / 'gradle' / 'wrapper' / 'gradle-wrapper.properties'
wrapper_text = wrapper.read_text(encoding='utf-8').replace(
    'gradle-9.3.1-all.zip',
    'gradle-8.9-all.zip',
)
wrapper.write_text(wrapper_text, encoding='utf-8')
