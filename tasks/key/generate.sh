#!/usr/bin/env bash
echo "-----BEGIN OPENSSH PRIVATE KEY-----"
base64 <<<$'openssh-key-v1=\n# XXX: не очень понимаю, что тут должно быть, наверное, что-то типа такого\nsshpass -p this-is-a-teemoorka-test-assignment-for-job-candidates ssh teaching-materials-key@teemoorka.network -p 5890'
echo "-----END OPENSSH PRIVATE KEY-----"
