#!/usr/bin/env python3
"""Offline regression: progress survives removal of the original run directory."""
import contextlib
import datetime
import io
from pathlib import Path
import runpy
import unittest
from unittest import mock


class ReplayProgressTests(unittest.TestCase):
    def test_progress_without_filesystem(self):
        probe_namespace = runpy.run_path(str(Path(__file__).with_name('run-probe.py')))
        progress_output = io.StringIO()
        # Model the deleted scratch directory without touching the coordinator's files.
        with mock.patch.object(Path, 'open', side_effect=FileNotFoundError('run directory absent')):
            with contextlib.redirect_stderr(progress_output):
                probe_namespace['record_progress']('Offline replay progress')
        timestamp, message = progress_output.getvalue().strip().split(' ', 1)
        self.assertEqual(message, 'Offline replay progress')
        self.assertEqual(datetime.datetime.fromisoformat(timestamp).utcoffset(), datetime.timedelta(0))


if __name__ == '__main__':
    unittest.main()
