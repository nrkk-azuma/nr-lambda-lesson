import json
import boto3
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
       
client = boto3.client('lambda')

def lambda_handler(event, context):
    inputParams = {
        'message': 'Hello world!!',
        'executor': 'PythonExecutor'
    }
    response = client.invoke(
        FunctionName = 'arn:aws:lambda:ap-northeast-1:**********:function:PythonTest',
        InvocationType = 'RequestResponse',
        Payload = json.dumps(inputParams)
    )
 
    responseFromChild = json.load(response['Payload'])
    if 'errorMessage' in responseFromChild:
      root_logger.error(responseFromChild);
      return {
        'statusCode': 500,
        'body': json.dumps('PythonTest was failed.')
      }
    else :
      root_logger.info(responseFromChild['body'])
      return {
        'statusCode': 200,
        'body': json.dumps('Execute PythonTest succeeded!')
      }

