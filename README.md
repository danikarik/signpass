# Signpass

CLI for Apple PassKit Generation

## Instructions

1. Create directory where all required certificates will be stored.
    ```sh
    mkdir -p certificates
    ```

2. Place `Apple Worldwide Developer Relations Certification Authority` inside `certificates` directory.

3. Place your `private` certificate (.p12) inside directory.

4. Extract pass certificate in `PEM` format:
    ```sh
    signpass cert -i ./certificates/<your_private_cert>.p12 -o ./certificates/passcertificate.pem -p <password_used_in_export>

5. Extract pass key in `PEM` format:
    ```sh
    signpass key -i ./certificates/<your_private_cert>.p12 -o ./certificates/passkey.pem -p <password_used_in_export>

6. Generate `.pkpass` archive:
    ```sh
    signpass -w ./certificates/<WWDR>.pem -s ./certificates/passcertificate.pem -k ./certificates/passkey.pem -r <raw_package> -p secret -d <output_dir>
    ```
