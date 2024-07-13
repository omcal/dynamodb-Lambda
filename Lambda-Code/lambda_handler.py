import json
import boto3

def lambda_handler(event, context):
    client_dynamo = boto3.resource('dynamodb')
    table = client_dynamo.Table('picus')

    results = []
    
    try:
        payload = event
        print(payload)
        
        if isinstance(payload, list):  
            for key_obj in payload:
                if isinstance(key_obj, dict) and 'id' in key_obj:  
                    key = key_obj['id']
                    try:
                        table.delete_item(Key={'id': key})
                        results.append({'id': key, 'status': 'success'})
                    except Exception as e:
                        results.append({'id': key, 'status': 'error', 'message': str(e)})
    except Exception as e:
        return {
            'statusCode': 400,
            'body': json.dumps({'error': 'Invalid JSON payload', 'message': str(e)})
        }

    return {
        'statusCode': 200,
        'body': json.dumps(results)
    }

