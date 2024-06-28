from flask import Flask, request
import logging

# Create the Flask application
app = Flask(__name__)

# Configure logging
logging.basicConfig(filename='error.log', level=logging.INFO)

@app.route('/500', methods=['GET', 'POST'])
def trigger_500():
    # Log the request
    logging.info(f'Request received: {request.method} {request.path} {request.remote_addr}')
    
    # Return a 500 status code
    return 'Internal Server Error', 500

if __name__ == '__main__':
    # Run the Flask application
    app.run(debug=True)
