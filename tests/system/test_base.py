from httppollerbeat import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting Httppollerbeat normally
        """
        self.render_config_template(
            path=os.path.abspath(self.working_dir) + "/log/*"
        )

        httppollerbeat_proc = self.start_beat()
        self.wait_until(lambda: self.log_contains("httppollerbeat is running"))
        exit_code = httppollerbeat_proc.kill_and_wait()
        assert exit_code == 0
