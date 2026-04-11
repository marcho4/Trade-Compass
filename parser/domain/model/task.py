
from dataclasses import dataclass
from enum import Enum

class TaskType(str, Enum):
    EXPECT_RAW_DATA = "expect-raw-data"
    EXTRACT = "extract"

@dataclass
class Task:
    id: str
    ticker: str
    year: int
    period: str
    report_url: str
    type: TaskType