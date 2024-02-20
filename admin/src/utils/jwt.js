import jwt from '@tsndr/cloudflare-worker-jwt'

export async function gen(data) {
    const token = await jwt.sign({
        ...data,
        nbf: Math.floor(Date.now() / 1000), // Not before: Now
        exp: Math.floor(Date.now() / 1000) + (10 * (60 * 60)) // Expires: Now + 2h
    }, 'chatgpt')

    return token
}


export async function verify(token) {
    // Verifing token
    try {
        const isValid = await jwt.verify(token, 'chatgpt', {throwError: true}) // false
        // Check for validity
        if (!isValid)
            return

        // Decoding token
        const { payload } = jwt.decode(token)
        return payload
    } catch (e) {
        console.log(e)
    }
}