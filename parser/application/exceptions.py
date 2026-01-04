class ParserError(Exception):
    pass


class CompanyNotFoundError(ParserError):
    pass


class DownloadError(ParserError):
    pass


class S3UploadError(ParserError):
    pass


class ConfigurationError(ParserError):
    pass


class PeriodParseError(ParserError):
    pass
