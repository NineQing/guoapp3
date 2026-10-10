class PickedFile {
  Future<int?> length() async => null;

  Future<List<int>> readAsBytes() async => const [];
}

enum FileType { any, audio, custom, image, media, video }

class FilePicker {
  static Future<String?> saveFile({
    String? fileName,
    List<int>? bytes,
    String? mimeType,
  }) async => null;

  static Future<PickedFile?> pickFile({
    FileType? type,
    List<String>? allowedExtensions,
  }) async => null;

  static Future<String?> getDirectoryPath({String? dialogTitle}) async => null;
}
