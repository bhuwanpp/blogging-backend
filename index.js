const express = require('express')
const bodyParser = require('body-parser')
const cors = require('cors')
const app = express()
app.use(bodyParser.json())
const { Client, Pool } = require('pg')
const port = 4000
app.use(cors())


const client = new Client({
    connectionString: 'postgres://root:GzWPeJv6HODIqI1mjH8UyvvprIGc3rAv@dpg-co6irf8l6cac73a9p1p0-a.singapore-postgres.render.com/blogging_kfua',
    ssl: {
        rejectUnauthorized: false, // This is required to connect to Render's PostgreSQL instance
    },
});

client
    .connect()
    .then(() => {
        console.log('Connected to postgres database')

        client.query('SELECT * FROM allblogs', (err, result) => {
            if (err) {
                console.log('Error executing query', err)

            }
            else {
                console.log('Query result ', result.rows)
            }

            client
                .end()
                .then(() => {
                    console.log('Connection to PostgreSQL closed');
                })
                .catch((err) => {
                    console.error('Error closing connection', err);
                });
        })
    })

    .catch((err) => {
        console.log('Error connecting to postgres database', err)
    })


const pool = new Pool({
    connectionString: 'postgres://root:GzWPeJv6HODIqI1mjH8UyvvprIGc3rAv@dpg-co6irf8l6cac73a9p1p0-a.singapore-postgres.render.com:5432/blogging_kfua',
    ssl: {
        rejectUnauthorized: false,
    },
});


app.get('/', (req, res) => {
    res.send('Backend for todo')
})
// post tasks
app.post('/blogs', async (req, res) => {
    try {
        const { blog } = req.body;
        const query = 'INSERT INTO allblogs (blog) VALUES($1) RETURNING *';
        const values = [blog]
        const result = await pool.query(query, values);
        res.json(result.rows[0]);

    } catch (err) {
        console.error(err)
        res.status(500).json({ err: 'an error occurred' })

    }
})

app.get('/get_blogs', async (req, res) => {
    try {
        const query = 'SELECT *  FROM allblogs'
        const result = await pool.query(query)
        res.json(result.rows);


    } catch (err) {
        console.log(err);
        res.status(500).json({ err: 'an error occurred', err })
    }
})
// Delete 
app.delete('/blogs/:id', async (req, res) => {
    try {
        const id = req.params.id
        const query = 'DELETE FROM allblogs WHERE id = $1'
        const result = await pool.query(query, [id])
        res.json({ message: 'Task delete successfully' })
    } catch (err) {
        console.log(err)
        res.status(500).json({ err: 'an error occurred' }, err)

    }

})
app.listen(port, () => {
    console.log('listening to the port', port)
})