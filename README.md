# Signpass

CLI for Apple PassKit Generation

## Development only

params: ["/Users/danikarik/code/kit/wallet/ServerReference/pass_server/lib/../data/passes/sample", "/Users/danikarik/code/kit/wallet/ServerReference/pass_server/lib/../data/Certificate/danikarik.p12", "dthcnf07", "/Users/danikarik/code/kit/wallet/ServerReference/pass_server/lib/../data/Certificate/Apple Worldwide Developer Relations Certification Authority.pem", "/Users/danikarik/code/kit/wallet/ServerReference/pass_server/data/passes/sample.pkpass"]
Raw pass has manifest? false
Raw pass has signature? false
Creating temp dir at /var/folders/7c/hbcl1tpd4fs37qmdqbkymqy40000gn/T/d20190112-19071-1sn9d88
Copying pass to temp directory.
Cleaning .DS_Store files
Generating JSON manifest
Signing the manifest
Compressing the pass
Gate changed to 98.
Reference server setup complete.

openssl pkcs12 -in Certificates.p12 -clcerts -nokeys -out passcertificate.pem -passin pass:
openssl pkcs12 -in Certificates.p12 -nocerts -out passkey.pem -passin pass: -passout pass:12345
openssl smime -binary -sign -certfile WWDR.pem -signer passcertificate.pem -inkey passkey.pem -in manifest.json -out signature -outform DER -passin pass:12345
zip -r freehugcoupon.pkpass manifest.json pass.json signature logo.png logo@2x.png icon.png icon@2x.png strip.png strip@2x.png
