import secrets

def generate_bearer_token(length=32):
    """Generate a secure random Bearer token."""
    return secrets.token_hex(length)

def add_token_to_file(token, file_path='tokens.conf'):
    """Add the generated token to the tokens.conf file in the required format."""
    with open(file_path, 'a') as file:
        file.write(f"\"{token}\" 1;\n")

if __name__ == "__main__":
    num_tokens = int(input("Enter the number of tokens to generate: "))
    for _ in range(num_tokens):
        token = generate_bearer_token()
        add_token_to_file(token)
    print(f"Generated {num_tokens} token/s successfully.")