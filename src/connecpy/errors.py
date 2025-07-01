from enum import Enum


class Errors(Enum):
    """
    Enum class representing different error codes and their corresponding status codes.
    """

    Canceled = "canceled"
    Unknown = "unknown"
    InvalidArgument = "invalid_argument"
    DeadlineExceeded = "deadline_exceeded"
    NotFound = "not_found"
    AlreadyExists = "already_exists"
    PermissionDenied = "permission_denied"
    Unauthenticated = "unauthenticated"
    ResourceExhausted = "resource_exhausted"
    FailedPrecondition = "failed_precondition"
    Aborted = "aborted"
    OutOfRange = "out_of_range"
    Unimplemented = "unimplemented"
    Internal = "internal"
    Unavailable = "unavailable"
    DataLoss = "data_loss"

    # Errors not defined in connect protocol
    NoError = ""
    BadRoute = "bad_route"
    Malformed = "malformed"

    @staticmethod
    def get_status_code(code: "Errors") -> int:
        """
        Deprecated: Use `to_http_status` instead.

        Returns the corresponding HTTP status code for the given error code.

        Args:
            code (Errors): The error code.

        Returns:
            int: The corresponding HTTP status code.
        """
        return code.to_http_status()

    def to_http_status(self) -> int:
        """
        Returns the corresponding HTTP status code for this error code.

        Returns:
            int: The HTTP status code.
        """
        return _error_to_http_status.get(self, 500)

    @staticmethod
    def from_http_status(code: int) -> "Errors":
        """
        Returns the corresponding HTTP status code for the given error code.

        Args:
            int: The HTTP status code.

        Returns:
            Errors: The corresponding error code.
        """
        return _http_status_to_error.get(code, Errors.Unknown)


_error_to_http_status = {
    Errors.Canceled: 499,
    Errors.Unknown: 500,
    Errors.InvalidArgument: 400,
    Errors.DeadlineExceeded: 504,
    Errors.NotFound: 404,
    Errors.AlreadyExists: 409,
    Errors.PermissionDenied: 403,
    Errors.ResourceExhausted: 429,
    Errors.FailedPrecondition: 400,
    Errors.Aborted: 409,
    Errors.OutOfRange: 400,
    Errors.Unimplemented: 501,
    Errors.Internal: 500,
    Errors.Unavailable: 503,
    Errors.DataLoss: 500,
    Errors.Unauthenticated: 401,
    # Errors not defined in connect protocol
    Errors.NoError: 200,
    Errors.BadRoute: 404,
    Errors.Malformed: 400,
}

_http_status_to_error = {
    400: Errors.Internal,
    401: Errors.Unauthenticated,
    403: Errors.PermissionDenied,
    404: Errors.Unimplemented,
    429: Errors.Unavailable,
    502: Errors.Unavailable,
    503: Errors.Unavailable,
    504: Errors.Unavailable,
}
