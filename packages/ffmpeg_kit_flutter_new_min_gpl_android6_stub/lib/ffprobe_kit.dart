class MediaInformation {
  Map<String, dynamic> getAllProperties() => const {};
}

class FFprobeSession {
  MediaInformation? getMediaInformation() => null;
}

class FFprobeKit {
  static Future<FFprobeSession> getMediaInformationFromCommandArguments(
    List<String> arguments,
  ) async => FFprobeSession();
}
