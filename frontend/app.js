const bucketName = 'thingies-input'
const region = 'eu-west-1'
const identityPool = 'eu-west-1:c4de6544-663a-40bd-a017-2e0924f55ca1'

AWS.config.update({
    region: region,
    credentials: new AWS.CognitoIdentityCredentials({
        IdentityPoolId: identityPool
    })
})

const s3 = new AWS.S3({
    apiVersion: '2006-03-01',
    params: { Bucket: bucketName }
})

function convert() {
    const files = document.getElementById("select-files").files
    if (!files.length) {
        msg('please select an image')
        return
    }
    const file = files[0]

    s3.upload({
        Key: generateToken(file.name),
        Body: file,
    }, function(err, data) {
        if (err) {
            msg('upload failed: ' + err.message)
            return
        }
        msg('upload sucessful')
    })
}

function generateToken(filename) {
    const ext = filename.substr(filename.lastIndexOf('.'))
    return rand() + ext
}

function msg(text) {
    document.getElementById("message").innerHTML = text
}

function rand() {
    const possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var text = ""
    for (var i = 0; i < 16; i++) {
        text += possible.charAt(Math.floor(Math.random() * possible.length))
    }
    return text
}