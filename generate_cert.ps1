# Get the directory where the script is running
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# Define file paths
$crtPath = Join-Path -Path $scriptDir -ChildPath "certificate.crt"
$keyPath = Join-Path -Path $scriptDir -ChildPath "certificate.key"

# Create the self-signed certificate
$cert = New-SelfSignedCertificate -DnsName "localhost" -CertStoreLocation "Cert:\CurrentUser\My" -KeyExportPolicy Exportable

# Export the certificate (.crt) in PEM format
$crtContent = "-----BEGIN CERTIFICATE-----`n" + [Convert]::ToBase64String($cert.RawData, 'InsertLineBreaks') + "`n-----END CERTIFICATE-----"
[System.IO.File]::WriteAllText($crtPath, $crtContent)

# Export the private key (.key) in PKCS#1 PEM format
$privateKey = [System.Security.Cryptography.X509Certificates.RSACertificateExtensions]::GetRSAPrivateKey($cert)
$keyBytes = $privateKey.ExportRSAPrivateKey()

$keyPem = "-----BEGIN RSA PRIVATE KEY-----`n" + [Convert]::ToBase64String($keyBytes, 'InsertLineBreaks') + "`n-----END RSA PRIVATE KEY-----"
[System.IO.File]::WriteAllText($keyPath, $keyPem)

Write-Output "Certificate created: $crtPath"
Write-Output "Private key saved: $keyPath"