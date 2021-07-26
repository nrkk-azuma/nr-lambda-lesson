import json
import newrelic.agent
import logging
from newrelic.agent import NewRelicContextFormatter
newrelic.agent.initialize()

handler = logging.StreamHandler(sys.stdout)
formatter = NewRelicContextFormatter()
handler.setFormatter(formatter)
root_logger = logging.getLogger()
for h in root_logger.handlers:
  root_logger.removeHandler(h)
root_logger.addHandler(handler)
root_logger.setLevel(logging.INFO)

@newrelic.agent.lambda_handler()
def lambda_handler(event, context):
    logger.info('This is Sample logs')
    return {
        'statusCode': 200,
        'body': json.dumps('Hello from Lambda!')
    }

